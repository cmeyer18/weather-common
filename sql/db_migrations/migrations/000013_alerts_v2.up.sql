CREATE TABLE alertV2 (
    id TEXT PRIMARY KEY,
    type TEXT,
    geometry geometry,
    areaDesc TEXT,
    sent TIMESTAMP WITH TIME ZONE,
    effective TIMESTAMP WITH TIME ZONE,
    onset TIMESTAMP WITH TIME ZONE,
    expires TIMESTAMP WITH TIME ZONE,
    ends TIMESTAMP WITH TIME ZONE,
    status TEXT,
    messageType TEXT,
    category TEXT,
    severity TEXT,
    certainty TEXT,
    urgency TEXT,
    event TEXT,
    sender TEXT,
    senderName TEXT,
    headline TEXT,
    description TEXT,
    instruction TEXT,
    response TEXT,
    parameters JSONB
);

CREATE TABLE alertV2_References (
    alertId TEXT,
    referenceId TEXT,
    FOREIGN KEY (alertId) REFERENCES alertV2(id)
);

CREATE TABLE alertV2_SAMECodes (
    alertId TEXT,
    code VARCHAR(20),
    FOREIGN KEY (alertId) REFERENCES alertV2(id)
);

CREATE TABLE alertV2_UGCCodes (
    alertId TEXT,
    code VARCHAR(20),
    FOREIGN KEY (alertId) REFERENCES alertV2(id)
);

CREATE TABLE convectiveOutlookV2 (
    outlookType TEXT,
    geometry geometry,
    dn int,
    issued TIMESTAMP WITH TIME ZONE,
    expires TIMESTAMP WITH TIME ZONE,
    valid TIMESTAMP WITH TIME ZONE,
    label TEXT,
    label2 TEXT,
    stroke TEXT,
    fill TEXT,
    PRIMARY KEY (outlookType, issued, label)
);

CREATE TABLE mesoscaleDiscussionV2 (
    number int,
    year int,
    geometry geometry,
    rawText text,
    PRIMARY KEY (number,year)
);
