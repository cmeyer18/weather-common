CREATE TABLE IF NOT EXISTS mesoscaleDiscussion(
    mdNumber int,
    year int,
    affectedArea jsonb,
    rawText text,
    CONSTRAINT pk_year_mdNumber PRIMARY KEY (mdNumber,year)
)