ALTER TABLE userNotification ADD isDeviceLocation BOOLEAN;

UPDATE userNotification
SET isDeviceLocation = FALSE
WHERE isDeviceLocation IS NULL
