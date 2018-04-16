package main

import (
	"fmt"
	"encoding/json"
	"github.com/oipwg/media-protocol/oip042"
	"strings"
	"time"
	"strconv"
	"os"
)

type OipArtifact struct {
	Pt oip042.PublishTomogram `json:"artifact"`
}
type OipPublish struct {
	OipArtifact `json:"publish"`
}
type rWrap struct {
	OipPublish `json:"oip042"`
}

func main() {
	ids, err := GetFilterIdList()
	if err != nil {
		panic(err)
	}

	for _, id := range ids {
		if _, ok := history[id]; ok {
			fmt.Printf("Tilt %s already published\n", id)
			continue
		}

		s, err := ipfsPinPath("/services/tomography/data/" + id)
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
		break
		pt, err := tiltIdToPublishTomogram(id)
		if err != nil {
			fmt.Println("Unable to obtain " + id)
			fmt.Println(err)
		} else {
			fmt.Println("---------")
			//PrettyPrint(pt)

			min, err := json.Marshal(rWrap{OipPublish{OipArtifact{pt}}})
			if err != nil {
				panic(err)
			}
			ids, err := sendToBlockchain("json:" + string(min))
			if err != nil {
				fmt.Println(ids)
				panic(err)
			} else {
				history[id] = ids
				PrettyPrint(ids)
			}
		}

		err = saveHistory()
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func processFiles(row TiltSeries) (ipfsHash, error) {
	return ipfsHash{}, nil
}

func tiltIdToPublishTomogram(tiltSeriesId string) (oip042.PublishTomogram, error) {
	tsr, err := GetTiltSeriesById(tiltSeriesId)
	if err != nil {
		panic(err)
	}

	//PrettyPrint(tsr)
	var pt oip042.PublishTomogram

	hash, ok := ipfsHashes[tiltSeriesId]
	if !ok || hash.Data == "" || hash.KeyMov == "" || hash.Combined == "" {
		hash, err = processFiles(tsr)
		if err != nil {
			return pt, err
		}
	}

	ts := time.Now().Unix()
	floAddress := config.FloAddress

	v := []string{hash.Combined, floAddress, strconv.FormatInt(ts, 10)}
	preImage := strings.Join(v, "-")
	signature, err := signMessage(floAddress, preImage)
	if err != nil {
		return pt, err
	}

	pt = oip042.PublishTomogram{
		PublishArtifact: oip042.PublishArtifact{
			Type:       "research",
			SubType:    "tomogram",
			Timestamp:  ts,
			FloAddress: floAddress,
			Signature:  signature,
			Info: &oip042.ArtifactInfo{
				Title:       tsr.Title,
				Description: "Auto imported from etdb",
				Tags:        "etdb,jensen.lab,tomogram,electron.tomography",
			},
			Storage: &oip042.ArtifactStorage{
				Network:  "ipfs",
				Location: hash.Combined,
				Files:    []oip042.ArtifactFiles{},
			},
			Payment: nil, // it's free
		},
		TomogramDetails: oip042.TomogramDetails{
			Microscopist:   tsr.Microscopist,
			Institution:    "Caltech",
			Lab:            "Jensen Lab",
			Magnification:  tsr.Magnification,
			Defocus:        tsr.Defocus,
			Strain:         tsr.SpeciesStrain,
			SpeciesName:    tsr.SpeciesName,
			ScopeName:      tsr.ScopeName,
			Date:           tsr.Date.Unix(),
			Emdb:           tsr.Emdb,
			TiltSingleDual: tsr.SingleDual,
			NBCItaxID:      tsr.SpeciesTaxId,
		},
	}

	if len(tsr.ScopeNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Scope notes: " + tsr.ScopeNotes + "\n"
	}
	if len(tsr.SpeciesNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Species notes: " + tsr.SpeciesNotes + "\n"
	}
	if len(tsr.TiltSeriesNotes) != 0 {
		pt.TomogramDetails.ArtNotes += "Tilt series notes: " + tsr.TiltSeriesNotes + "\n"
	}

	for _, df := range tsr.DataFiles {
		fi, err := os.Stat(df.FilePath)
		if err != nil {
			return pt, err
		}
		af := oip042.ArtifactFiles{
			Type:    df.Type,
			SubType: df.SubType,
			FNotes:  df.Notes,
			Fsize:   fi.Size(),
			Dname:   df.Filename,
			Fname:   strings.TrimPrefix(df.FilePath, "/services/tomography/data/"+tsr.Id+"/"),
		}
		pt.Storage.Files = append(pt.Storage.Files, af)
	}

	for _, tdf := range tsr.ThreeDFiles {
		fi, err := os.Stat(tdf.FilePath)
		if err != nil {
			return pt, err
		}
		af := oip042.ArtifactFiles{
			Type:     tdf.Type,
			SubType:  tdf.SubType,
			FNotes:   tdf.Notes,
			Fsize:    fi.Size(),
			Dname:    tdf.Filename,
			Fname:    strings.TrimPrefix(tdf.FilePath, "/services/tomography/data/"+tsr.Id+"/"),
			Software: tdf.Software,
		}
		pt.Storage.Files = append(pt.Storage.Files, af)
	}

	if tsr.KeyImg > 0 && tsr.KeyImg <= 4 {
		kif := "keyimg_" + tsr.Id + "_s.jpg"
		fi, err := os.Stat("/services/tomography/data/" + tsr.Id + "/" + kif)
		if err != nil {
			return pt, err
		}
		ki := oip042.ArtifactFiles{
			Type:    "image",
			SubType: "thumbnail",
			CType:   "image/jpeg",
			Fsize:   fi.Size(),
			Fname:   kif,
		}
		pt.Storage.Files = append(pt.Storage.Files, ki)

		kif = "keyimg_" + tsr.Id + ".jpg"
		fi, err = os.Stat("/services/tomography/data/" + tsr.Id + "/" + kif)
		if err != nil {
			return pt, err
		}
		ki = oip042.ArtifactFiles{
			Type:    "research",
			SubType: "keyimg",
			CType:   "image/jpeg",
			Fsize:   fi.Size(),
			Fname:   "keyimg_" + tsr.Id + ".jpg",
		}
		pt.Storage.Files = append(pt.Storage.Files, ki)
	}
	if tsr.KeyMov > 0 && tsr.KeyMov <= 4 {
		kmf := "keymov_" + tsr.Id + ".flv"
		fi, err := os.Stat("/services/tomography/data/" + tsr.Id + "/" + kmf)
		if err != nil {
			return pt, err
		}
		km := oip042.ArtifactFiles{
			Type:    "research",
			SubType: "keymov",
			CType:   "video/x-flv",
			Fsize:   fi.Size(),
			Fname:   kmf,
		}
		pt.Storage.Files = append(pt.Storage.Files, km)
	}

	return pt, nil
}

func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
