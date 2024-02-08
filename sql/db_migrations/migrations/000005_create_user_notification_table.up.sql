CREATE TABLE IF NOT EXISTS userNotification(
    notificationId       varchar(255) primary key,
    userId				 varchar(255),
    zoneCode			 varchar(255),
    countyCode			 varchar(255),
    creationTime		 timestamptz,
    lat					 float,
    lng					 float,
    formattedAddress	 varchar(500),
    apnKey				 varchar(255),
    locationName		 varchar(255)
)
