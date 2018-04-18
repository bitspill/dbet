-- Multi row list of Data Files
SELECT
  df.DEF_id,
  df.filetype,
  df.filename,
  df.TXT_notes,
  df.auto,
  df.`REF|ThreeDFile|image` as ThreeDFileImage
FROM DataFile as df  where `REF|TiltSeriesData|tiltseries` = ?;