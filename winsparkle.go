// +build windows

// Package winsparkle provides go bindings for WinSparkle.
//
// WinSparkle is a plug-and-forget software update library for Windows
// applications. It is heavily inspired by the Sparkle framework for OS X
// written by Andy Matuschak and others, to the point of sharing the same
// updates format (appcasts) and having a very similar user interface.
//
// See https://winsparkle.org for more information about WinSparkle.
package winsparkle

import (
	"errors"
	"syscall"
	"time"
)

var winsparkle = syscall.MustLoadDLL("WinSparkle.dll")

// Starts WinSparkle.
//
// If WinSparkle is configured to check for updates on startup, proceeds
// to perform the check. You should only call this function when your app
// is initialized and shows its main window.
//
// This call doesn't block and returns almost immediately. If an
// update is available, the respective UI is shown later from a separate
// thread.
func Init() {
	winsparkle.MustFindProc("win_sparkle_init").Call()
}

// Cleans up after WinSparkle.
//
// Should be called by the app when it's shutting down. Cancels any
// pending Sparkle operations and shuts down its helper threads.
func Cleanup() {
	winsparkle.MustFindProc("win_sparkle_cleanup").Call()
}

// Sets UI language from its ISO code.
//
// This function must be called before sparkle.Init().
//
// Param lang must be an ISO 639 language code with an optional ISO 3116 country
// code, e.g. "fr", "pt-PT", "pt-BR" or "pt_BR", as used
// e.g. by ::GetThreadPreferredUILanguages() too.
func SetLang(lang string) {
	winsparkle.MustFindProc("win_sparkle_set_lang").Call(string2uintptr(lang))
}

// Sets UI language from its Win32 LANGID code.
//
// This function must be called before win_sparkle_init().
//
// Param langid must be a Language code (LANGID) as created by the MAKELANGID macro
// or returned by e.g. ::GetThreadUILanguage()
func SetLangID(langid uint32) {
	winsparkle.MustFindProc("win_sparkle_set_langid").Call(uintptr(langid))
}

// Sets URL for the app's appcast.
//
// Only http and https schemes are supported.
//
// If this function isn't called by the app, the URL is obtained from Windows resource named "FeedURL" of type "APPCAST".
//
// Note: Always use HTTPS feeds, do not use unencrypted HTTP! This is
// necessary to prevent both leaking user information and preventing
// various MITM attacks.
//
// Note: See https://github.com/vslavik/winsparkle/wiki/Appcast-Feeds for
// more information about appcast feeds.
func SetAppcastURL(url string) {
	winsparkle.MustFindProc("win_sparkle_set_appcast_url").Call(string2uintptr(url))
}

// Sets DSA public key.
//
// Only PEM format is supported.
//
// Public key will be used to verify DSA signature of the update file.
// PEM data will be set only if it contains valid DSA public key.
//
// If this function isn't called by the app, public key is obtained from
// Windows resource named "DSAPub" of type "DSAPEM".
func SetDSAPubPEM(pem string) error {
	r, _, _ := winsparkle.MustFindProc("win_sparkle_set_dsa_pub_pem").Call(string2uintptr(pem))
	if r == 0 {
		return errors.New("invalid DSA public key provided")
	}
	return nil
}

// Sets application metadata.
//
// Normally, these are taken from VERSIONINFO/StringFileInfo resources,
// but if your application doesn't use them for some reason, using this
// function is an alternative.
//
// `app` is both shown to the user and used in HTTP User-Agent header.
//
// Note: `company` and `app` are used to determine the location of WinSparkle
// settings in registry (HKCU\Software\<company>\<app>\WinSparkle is used).
func SetAppDetails(company, appName, version string) {
	winsparkle.MustFindProc("win_sparkle_set_app_details").Call(
		string2wchar(company), string2wchar(appName), string2wchar(version))
}

// Sets application build version number.
//
// This is the internal version number that is not normally shown to the user.
// It can be used for finer granularity that official release versions, e.g. for
// interim builds.
//
// If this function is called, then the provided *build* number is used for comparing
// versions; it is compared to the "version" attribute in the appcast and corresponds
// to OS X Sparkle's CFBundleVersion handling. If used, then the appcast must
// also contain the "shortVersionString" attribute with human-readable display
// version string. The version passed to sparkle.SetAppDetails()
// corresponds to this and is used for display.
func SetAppBuildVersion(build string) {
	winsparkle.MustFindProc("win_sparkle_set_app_build_version").Call(string2wchar(build))
}

// Set custom HTTP header for appcast checks.
func SetHTTPHeader(name, value string) {
	winsparkle.MustFindProc("win_sparkle_set_http_header").Call(string2uintptr(name), string2uintptr(value))
}

// Clears all custom HTTP headers previously added using sparkle.SetHTTPHeader()
func ClearHTTPHeaders() {
	winsparkle.MustFindProc("win_sparkle_clear_http_headers").Call()
}

// Set the registry path where settings will be stored.
//
// Normally, these are stored in
// "HKCU\Software\<company_name>\<app_name>\WinSparkle"
// but if your application needs to store the data elsewhere for
// some reason, using this function is an alternative.
//
// Note that `path` is relative to HKCU/HKLM root and the root is not part
// of it. For example:
//
//	sparkle.SetRegistryPath("Software\\My App\\Updates");
func SetRegistryPath(path string) {
	winsparkle.MustFindProc("win_sparkle_set_registry_path").Call(string2uintptr(path))
}

// Sets whether updates are checked automatically or only through a manual call.
// If disabled, sparkle.CheckUpdateWithUI() must be used explicitly.
func SetAutomaticCheckForUpdates(check bool) {
	winsparkle.MustFindProc("win_sparkle_set_automatic_check_for_updates").Call(bool2uintptr(check))
}

// Gets the automatic update checking state.
//
// Returns true if updates are set to be checked automatically, false otherwise.
//
// Note: Defaults to 0 when not yet configured (as happens on first start).
func GetAutomaticCheckForUpdates() bool {
	r, _, _ := winsparkle.MustFindProc("win_sparkle_get_automatic_check_for_updates").Call()
	return r == 1
}

// Sets the automatic update interval between checks for updates.
//
// Note: The minimum update interval is 1 hour.
func SetUpdateCheckInterval(interval time.Duration) {
	winsparkle.MustFindProc("win_sparkle_set_update_check_interval").Call(uintptr(interval.Seconds()))
}

// Gets the automatic update interval.
//
// Default value is one day.
func GetUpdateCheckInterval() time.Duration {
	r, _, _ := winsparkle.MustFindProc("win_sparkle_get_update_check_interval").Call()
	return time.Duration(r) * time.Second
}

// Gets the time for the last update check.
//
// Default value is -1, indicating that the update check has never run.
func GetLastCheckTime() time.Time {
	r, _, _ := winsparkle.MustFindProc("win_sparkle_get_last_check_time").Call()
	return time.Unix(int64(r), 0)
}

// Set callback to be called when the updater encounters an error.
func SetErrorCallback(cb func()) {
	fn := func() uintptr { cb(); return 0 }
	winsparkle.MustFindProc("win_sparkle_set_error_callback").Call(syscall.NewCallback(fn))
}

// Set callback for querying the application if it can be closed.
//
// This callback will be called to ask the host if it's ready to shut down,
// before attempting to launch the installer. The callback returns TRUE if
// the host application can be safely shut down or FALSE if not (e.g. because
// the user has unsaved documents).
func SetCanShutdownCallback(cb func() bool) {
	fn := func() uintptr { return bool2uintptr(cb()) }
	winsparkle.MustFindProc("win_sparkle_set_can_shutdown_callback").Call(syscall.NewCallback(fn))
}

// Set callback for shutting down the application.
//
// This callback will be called to ask the host to shut down immediately after
// launching the installer. Its implementation should gracefully terminate the
// application.
func SetShutdownRequestCallback(cb func()) {
	fn := func() uintptr { cb(); return 0 }
	winsparkle.MustFindProc("win_sparkle_set_shutdown_request_callback").Call(syscall.NewCallback(fn))
}

// Set callback to be called when the updater did find an update.
//
// This is useful in combination with sparkle.CheckUpdateWithUIAndInstall()
// as it allows you to perform some action after WinSparkle checks for updates.
func SetDidFindUpdateCallback(cb func()) {
	fn := func() uintptr { cb(); return 0 }
	winsparkle.MustFindProc("win_sparkle_set_did_find_update_callback").Call(syscall.NewCallback(fn))
}

// Set callback to be called when the updater did not find an update.
//
// This is useful in combination with sparkle.CheckUpdateWithUIAndInstall()
// as it allows you to perform some action after WinSparkle checks for updates.
func SetDidNotFindUpdateCallback(cb func()) {
	fn := func() uintptr { cb(); return 0 }
	winsparkle.MustFindProc("win_sparkle_set_did_not_find_update_callback").Call(syscall.NewCallback(fn))
}

// Set callback to be called when the user cancels a download.
//
// This is useful in combination with sparkle.CheckUpdateWithUIAndInstall()
// as it allows you to perform some action when the installation is interrupted.
func SetUpdateCancelledCallback(cb func()) {
	fn := func() uintptr { cb(); return 0 }
	winsparkle.MustFindProc("win_sparkle_set_update_cancelled_callback").Call(syscall.NewCallback(fn))
}

// Checks if an update is available, showing progress UI to the user.
//
// Normally, WinSparkle checks for updates on startup and only shows its UI
// when it finds an update. If the application disables this behavior, it
// can hook this function to "Check for updates..." menu item.
//
// When called, background thread is started to check for updates. A small
// window is shown to let the user know the progress. If no update is found,
// the user is told so. If there is an update, the usual "update available"
// window is shown.
//
// This function returns immediately.
//
// Note: Because this function is intended for manual, user-initiated checks
// for updates, it ignores "Skip this version" even if the user checked it
// previously.
func CheckUpdateWithUI() {
	winsparkle.MustFindProc("win_sparkle_check_update_with_ui").Call()
}

// Checks if an update is available, showing progress UI to the user and
// immediately installing the update if one is available.
//
// This is useful for the case when users should almost always use the
// newest version of your software. When called, WinSparkle will check for
// updates showing a progress UI to the user. If an update is found the update
// prompt will be skipped and the update will be installed immediately.
//
// If your application expects to do something after checking for updates you
// may wish to use sparkle.SetDidNotFindUpdateCallback() and
// sparkle.SetUpdateCancelledCallback().
func CheckUpdateWithUIAndInstall() {
	winsparkle.MustFindProc("win_sparkle_check_update_with_ui_and_install").Call()
}

// Checks if an update is available.
//
// No progress UI is shown to the user when checking. If an update is
// available, the usual "update available" window is shown; this function
// is *not* completely UI-less.
//
// Use with caution, it usually makes more sense to use the automatic update
// checks on interval option or manual check with visible UI.
//
// This function returns immediately.
//
// Note: This function respects "Skip this version" choice by the user.
func CheckUpdateWithoutUI() {
	winsparkle.MustFindProc("win_sparkle_check_update_without_ui").Call()
}
