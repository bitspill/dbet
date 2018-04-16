-- Multi row list of Three D Files
SELECT
  tdf.TXT_notes,
  tdf.classify,
  tdf.filename,
  tdf.DEF_id
FROM ThreeDFile as tdf
WHERE tdf.`REF|TiltSeriesData|tiltseries` = ?;