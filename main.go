package notif

import (
	"runtime"

	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

// func main() {
// 	scheduleNotification("com.example.fake2", "unique-adsfabd",
// 		"Reminder", "You have have a task due soon", 4)

// 	removeScheduledNotification("unique-adsfabd")
// }



// Schedules a user notification to the to the notification center
// Always use the same bundleIdentifier and change only if if you miss the allow notification
// notificationIdentifier should be a unique value, can be used to remove scheduled notification
// title and informativeText can be set as strings and set the deliverySecondsFromNow to schedule the
// notification delivery time
func scheduleNotification(bundleIdentifier string, notificationIdentifier string,
	title string, informativeText string, deliverySecondsFromNow float64) {
		runtime.LockOSThread()

		nsbundle := foundation.Bundle_MainBundle().Class()
		objc.ReplaceMethod(nsbundle, objc.Sel("bundleIdentifier"), func(_ objc.IObject) string {
			return bundleIdentifier // change this if you miss the allow notification
		})

		objc.WithAutoreleasePool(func() {
			notif := objc.Call[objc.Object](objc.GetClass("NSUserNotification"), objc.Sel("new"))
			notif.Autorelease()

			date := foundation.NewDate()
			date.Autorelease()

			var time foundation.TimeInterval = foundation.TimeInterval(deliverySecondsFromNow)
			date = date.InitWithTimeIntervalSinceNow(time)

			objc.Call[objc.Void](notif, objc.Sel("setTitle:"), title)
			objc.Call[objc.Void](notif, objc.Sel("setInformativeText:"), informativeText)
			objc.Call[objc.Void](notif, objc.Sel("setDeliveryDate:"), date)
			objc.Call[objc.Void](notif, objc.Sel("setIdentifier:"), notificationIdentifier)

			center := objc.Call[objc.Object](objc.GetClass("NSUserNotificationCenter"), objc.Sel("defaultUserNotificationCenter"))
			objc.Call[objc.Void](center, objc.Sel("scheduleNotification:"), notif)
			objc.Call[objc.Object](center, objc.Sel("scheduledNotifications"))
		})

		runtime.UnlockOSThread()
}

// Removes a scheduled notification from the queue. Needs the unique notificationIdentifier
// that was assigned to notificaiton when first scheduled
func removeScheduledNotification(notificationIdentifier string) {
		runtime.LockOSThread()


		objc.WithAutoreleasePool(func() {

			center := objc.Call[objc.Object](objc.GetClass("NSUserNotificationCenter"), objc.Sel("defaultUserNotificationCenter"))
			scheduledNotifications := objc.Call[objc.Object](center, objc.Sel("scheduledNotifications"))
			array := foundation.ArrayFrom(scheduledNotifications.Ptr())
			n := array.Count()

			for i := range n {
				nsNotif := array.ObjectAtIndex(i)
				notifIdentifier := objc.Call[objc.Object](nsNotif, objc.Sel("identifier"))
				str_identifier := foundation.StringFrom(notifIdentifier.Ptr())

				if str_identifier.IsEqualToString(notificationIdentifier) {
					objc.Call[objc.Void](center, objc.Sel("removeScheduledNotification:"), nsNotif)
				}
			}
			scheduledNotifications = objc.Call[objc.Object](center, objc.Sel("scheduledNotifications"))
		})

		runtime.UnlockOSThread()
}
