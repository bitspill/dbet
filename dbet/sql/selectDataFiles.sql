-- Multi row list of Data Files
SELECT
  df.filetype,
  df.filename,
  df.TXT_notes,
  df.`REF|ThreeDFile|image` as ThreeDFileImage
FROM DataFile as df  where `REF|TiltSeriesData|tiltseries` = ?;