UPDATE userNotification
    SET mesoscaleDiscussionNotifications = FALSE
    WHERE mesoscaleDiscussionNotifications IS NULL
