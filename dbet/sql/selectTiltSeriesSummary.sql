-- Single row summary data
SELECT
  tsd.tiltSeriesID,
  tsd.title,
  tsd.tomo_date,
  tsd.TXT_notes as tsd_TXT_notes,
  tsd.scope,
  tsd.roles,
  tsd.keyimg,
  tsd.keymov,
  scd.TXT_notes as scd_TXT_notes,
  spd.SpeciesName,
  spd.TXT_notes as spd_TXT_notes,
  spd.strain,
  spd.tax_id,
  tsd.single_dual,
  tsd.defocus,
  tsd.magnification,
  tsd.software_acquisition,
  tsd.software_process,
  tsd.emdb,
  u.fullname
FROM TiltSeriesData AS tsd LEFT JOIN ScopeData AS scd ON scd.scopename = tsd.scope
  JOIN SpeciesData AS spd ON tsd.`REF|SpeciesData|specie` = spd.DEF_id
  JOIN UserData AS u ON u.DEF_id = tsd.`REF|UserData|user`
WHERE ispublic = ? AND tiltseriesID = ? LIMIT 1;