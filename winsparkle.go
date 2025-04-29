//go:build windows
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

var winsparkle = syscall.NewLazyDLL("WinSparkle.dll")

// Init starts WinSparkle.
//
// If WinSparkle is configured to check for updates on startup, proceeds
// to perform the check. You should only call this function when your app
// is initialized and shows its main window.
//
// This call doesn't block and returns almost immediately. If an
// update is available, the respective UI is shown later from a separate
// thread.
func Init() {
	winsparkle.NewProc("win_sparkle_init").Call()
}

// Cleanup cleans up after WinSparkle.
//
// Should be called by the app when it's shutting down. Cancels any
// pending Sparkle operations and shuts down its helper threads.
func Cleanup() {
	winsparkle.NewProc("win_sparkle_cleanup").Call()
}

// SetLang sets UI language from its ISO code.
//
// This function must be called before [Init].
//
// Param lang must be an ISO 639 language code with an optional ISO 3116
// country code, e.g. "fr", "pt-PT", "pt-BR" or "pt_BR", as used
// e.g. by ::GetThreadPreferredUILanguages() too.
func SetLang(lang string) {
	winsparkle.NewProc("win_sparkle_set_lang").Call(char(lang))
}

// SetLangID sets UI language from its Win32 LANGID code.
//
// This function must be called before [Init].
//
// Param langid must be a Language code (LANGID) as created by the MAKELANGID
// macro or returned by e.g. ::GetThreadUILanguage().
//
// See https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid/
func SetLangID(langid uint16) {
	winsparkle.NewProc("win_sparkle_set_langid").Call(uintptr(langid))
}

// SetAppcastURL sets URL for the app's appcast.
//
// Only http and https schemes are supported.
//
// If this function isn't called by the app, the URL is obtained from Windows
// resource named "FeedURL" of type "APPCAST".
//
// Note: Always use HTTPS feeds, do not use unencrypted HTTP! This is
// necessary to prevent both leaking user information and preventing
// various MITM attacks.
//
// Note: See https://github.com/vslavik/winsparkle/wiki/Appcast-Feeds for
// more information about appcast feeds.
func SetAppcastURL(url string) {
	winsparkle.NewProc("win_sparkle_set_appcast_url").Call(char(url))
}

// SetDSAPubPEM sets DSA public key.
//
// Only PEM format is supported.
//
// Public key will be used to verify DSA signature of the update file.
// PEM data will be set only if it contains valid DSA public key.
//
// If this function isn't called by the app, public key is obtained from
// Windows resource named "DSAPub" of type "DSAPEM".
//
// Deprecated: DSA signatures are deprecated and will be removed in a future version.
// Migrate over to EdDSA (ed25519) using [SetEdDSAPublicKey], see
// https://github.com/vslavik/winsparkle/wiki/Upgrading-from-DSA-to-EdDSA-signatures.
func SetDSAPubPEM(pem string) error {
	r, _, _ := winsparkle.NewProc("win_sparkle_set_dsa_pub_pem").Call(char(pem))
	if r == 0 {
		return errors.New("invalid DSA public key provided")
	}
	return nil
}

// SetEdDSAPublicKey sets EdDSA public key.
//
// Only base64 encoded format is supported.
//
// Public key will be used to verify EdDSA signature of the update file.
// It will be set only if it contains valid EdDSA public key.
//
// If this function isn't called by the app, public key is obtained from
// Windows resource named "EdDSAPub" of type "EDDSA".
//
// Note: If this function is called, DSA public key set with [SetDSAPubPEM]
// or present in the resources will be ignored; so will DSA signatures in the appcast.
func SetEdDSAPublicKey(key string) error {
	r, _, _ := winsparkle.NewProc("win_sparkle_set_eddsa_public_key").Call(char(key))
	if r == 0 {
		return errors.New("invalid edDSA public key provided")
	}
	return nil
}

// SetAppDetails sets application metadata.
//
// Normally, these are taken from VERSIONINFO/StringFileInfo resources,
// but if your application doesn't use them for some reason, using this
// function is an alternative.
//
// `app` is both shown to the user and used in HTTP User-Agent header.
//
// Note: `company` and `app` are used to determine the location of WinSparkle
// settings in registry (HKCU\Software\<company>\<app>\WinSparkle is used).
func SetAppDetails(company, app, version string) {
	winsparkle.NewProc("win_sparkle_set_app_details").Call(wchar(company), wchar(app), wchar(version))
}

// SetAppBuildVersion sets application build version number.
//
// This is the internal version number that is not normally shown to the user.
// It can be used for finer granularity that official release versions, e.g. for
// interim builds.
//
// If this function is called, then the provided *build* number is used for
// comparing versions; it is compared to the "version" attribute in the appcast
// and corresponds to OS X Sparkle's CFBundleVersion handling. If used, then
// the appcast must also contain the "shortVersionString" attribute with
// human-readable display version string. The version passed to [SetAppDetails]
// corresponds to this and is used for display.
func SetAppBuildVersion(build string) {
	winsparkle.NewProc("win_sparkle_set_app_build_version").Call(wchar(build))
}

// SetHTTPHeader sets custom HTTP header for appcast checks.
func SetHTTPHeader(name, value string) {
	winsparkle.NewProc("win_sparkle_set_http_header").Call(char(name), char(value))
}

// ClearHTTPHeaders clears all custom HTTP headers previously added using
// [SetHTTPHeader].
func ClearHTTPHeaders() {
	winsparkle.NewProc("win_sparkle_clear_http_headers").Call()
}

// SetRegistryPath sets the registry path where settings will be stored.
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
	winsparkle.NewProc("win_sparkle_set_registry_path").Call(char(path))
}

// ConfigStore is used to override WinSparkle configuration's read, write and delete
// functions.
type ConfigStore interface {
	// Read returns a config value and a bool indicating if it was successful.
	Read(name string) (value string, ok bool)

	// Write a config value. Returns a bool indicating if it was successful.
	Write(name, value string) bool

	// Delete config value. Returns a bool indicating if it was successful.
	Delete(name string) bool
}

// SetConfigMethods overrides WinSparkle's configuration read, write and delete
// functions.
//
// By default, WinSparkle will read, write and delete configuration values by
// interacting directly with Windows Registry.
// If you want to manage configuration by yourself, or if you don't want let
// WinSparkle write settings directly to the Windows Registry, you can provide
// your own functions to read, write and delete configuration.
func SetConfigMethods(store ConfigStore) {
	winsparkle.NewProc("win_sparkle_set_config_methods").Call(uintptr(configMethods(store)))
}

// SetAutomaticCheckForUpdates sets whether updates are checked automatically
// or only through a manual call. If disabled, [CheckUpdateWithUI] must be used
// explicitly.
func SetAutomaticCheckForUpdates(check bool) {
	winsparkle.NewProc("win_sparkle_set_automatic_check_for_updates").Call(boolean(check))
}

// GetAutomaticCheckForUpdates gets the automatic update checking state.
//
// Returns true if updates are set to be checked automatically, false otherwise.
//
// Note: Defaults to 0 when not yet configured (as happens on first start).
func GetAutomaticCheckForUpdates() bool {
	r, _, _ := winsparkle.NewProc("win_sparkle_get_automatic_check_for_updates").Call()
	return r == 1
}

// SetUpdateCheckInterval sets the automatic update interval between checks for
// updates.
//
// Note: The minimum update interval is 1 hour.
func SetUpdateCheckInterval(interval time.Duration) {
	winsparkle.NewProc("win_sparkle_set_update_check_interval").Call(uintptr(interval / time.Second))
}

// GetUpdateCheckInterval gets the automatic update interval.
//
// Default value is one day.
func GetUpdateCheckInterval() time.Duration {
	r, _, _ := winsparkle.NewProc("win_sparkle_get_update_check_interval").Call()
	return time.Duration(r) * time.Second
}

// GetLastCheckTime gets the time for the last update check.
//
// Default value is -1, indicating that the update check has never run.
func GetLastCheckTime() time.Time {
	r, _, _ := winsparkle.NewProc("win_sparkle_get_last_check_time").Call()
	return time.Unix(int64(r), 0)
}

// SetErrorCallback sets callback to be called when the updater encounters an
// error.
func SetErrorCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_error_callback").Call(fn)
}

// SetCanShutdownCallback sets callback for querying the application if it can
// be closed.
//
// This callback will be called to ask the host if it's ready to shut down,
// before attempting to launch the installer. The callback returns `true` if
// the host application can be safely shut down or `false` if not
// (e.g. because the user has unsaved documents).
func SetCanShutdownCallback(cb func() bool) {
	fn := syscall.NewCallbackCDecl(func() uintptr { return boolean(cb()) })
	winsparkle.NewProc("win_sparkle_set_can_shutdown_callback").Call(fn)
}

// SetShutdownRequestCallback sets callback for shutting down the application.
//
// This callback will be called to ask the host to shut down immediately after
// launching the installer. Its implementation should gracefully terminate the
// application.
func SetShutdownRequestCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_shutdown_request_callback").Call(fn)
}

// SetDidFindUpdateCallback sets callback to be called when the updater did
// find an update.
//
// This is useful in combination with [CheckUpdateWithUIAndInstall]
// as it allows you to perform some action after WinSparkle checks for updates.
func SetDidFindUpdateCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_did_find_update_callback").Call(fn)
}

// SetDidNotFindUpdateCallback sets callback to be called when the updater did
// not find an update.
//
// This is useful in combination with [CheckUpdateWithUIAndInstall]
// as it allows you to perform some action after WinSparkle checks for updates.
func SetDidNotFindUpdateCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_did_not_find_update_callback").Call(fn)
}

// SetUpdateCancelledCallback sets callback to be called when the user cancels
// a download.
//
// This is useful in combination with [CheckUpdateWithUIAndInstall]
// as it allows you to perform some action when the installation is
// interrupted.
func SetUpdateCancelledCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_update_cancelled_callback").Call(fn)
}

// SetUpdateSkippedCallback sets callback to be called when the user skips an
// update.
//
// This is useful in combination with [CheckUpdateWithUIAndInstall]
// or similar as it allows you to perform some action when the update is
// skipped.
func SetUpdateSkippedCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_update_skipped_callback").Call(fn)
}

// SetUpdatePostponedCallback sets callback to be called when the user
// postpones an update (presses 'remind me later' button).
//
// This is useful in combination with [CheckUpdateWithUI] or
// similar as it allows you to perform some action when the download is
// postponed.
func SetUpdatePostponedCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_update_postponed_callback").Call(fn)
}

// SetUpdateDismissedCallback sets callback to be called when the user
// dismisses (closes) update dialog.
//
// This is useful in combination with [CheckUpdateWithoutUI] or similar
// as it allows you to perform some action when the update dialog is closed.
func SetUpdateDismissedCallback(cb func()) {
	fn := syscall.NewCallbackCDecl(func() uintptr { cb(); return 0 })
	winsparkle.NewProc("win_sparkle_set_update_dismissed_callback").Call(fn)
}

// SetUserRunInstallerCallback sets callback to be called when the update
// payload is downloaded and read to be executed or handled in some other
// manner.
//
// The callback returns a boolean indicating whether the update was handled
// and an error. If `handled` is `false` and there is no error WinSparkle's
// default handling will take place.
func SetUserRunInstallerCallback(cb func(file string) (handled bool, err error)) {
	fn := syscall.NewCallbackCDecl(func(p *uint16) int {
		ok, err := cb(utf16PtrToString(p))
		if err != nil {
			return -1
		}
		return int(boolean(ok))
	})
	winsparkle.NewProc("win_sparkle_set_user_run_installer_callback").Call(fn)
}

// CheckUpdateWithUI checks if an update is available, showing progress UI to
// the user.
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
	winsparkle.NewProc("win_sparkle_check_update_with_ui").Call()
}

// CheckUpdateWithUIAndInstall checks if an update is available, showing
// progress UI to the user and immediately installing the update if one is
// available.
//
// This is useful for the case when users should almost always use the
// newest version of your software. When called, WinSparkle will check for
// updates showing a progress UI to the user. If an update is found the update
// prompt will be skipped and the update will be installed immediately.
//
// If your application expects to do something after checking for updates you
// may wish to use [SetDidNotFindUpdateCallback] and
// [SetUpdateCancelledCallback].
func CheckUpdateWithUIAndInstall() {
	winsparkle.NewProc("win_sparkle_check_update_with_ui_and_install").Call()
}

// CheckUpdateWithoutUI checks if an update is available.
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
	winsparkle.NewProc("win_sparkle_check_update_without_ui").Call()
}
