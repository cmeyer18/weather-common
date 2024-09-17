CREATE TABLE IF NOT EXISTS location(
    locationID           text primary key,
    locationName		 text,
    locationType         SMALLINT,
    zoneCode			 text,
    countyCode			 text,
    created		         timestamp DEFAULT current_timestamp,
    latitude			 float,
    longitude			 float,
    locationReferenceID  text
);

CREATE TABLE IF NOT EXISTS locationOptions(
    locationID text,
    optionType SMALLINT,
    option text,

    PRIMARY KEY (locationID, optionType, option)
);
