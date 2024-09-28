ALTER TABLE userNotification ADD liveActivities BOOLEAN;

UPDATE userNotification
SET liveActivities = FALSE
WHERE liveActivities IS NULL
