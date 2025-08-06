//go:build darwin
// +build darwin

package overlaydetector

/*
#cgo LDFLAGS: -framework CoreGraphics
#include <CoreGraphics/CoreGraphics.h>
#include <CoreFoundation/CoreFoundation.h>

static int getLayer(CFDictionaryRef dict) {
    CFNumberRef num = CFDictionaryGetValue(dict, kCGWindowLayer);
    int layer = 0;
    if (num) CFNumberGetValue(num, kCFNumberIntType, &layer);
    return layer;
}
static double getAlpha(CFDictionaryRef dict) {
    CFNumberRef num = CFDictionaryGetValue(dict, kCGWindowAlpha);
    double alpha = 1.0;
    if (num) CFNumberGetValue(num, kCFNumberFloat64Type, &alpha);
    return alpha;
}
static CGRect getBounds(CFDictionaryRef dict) {
    CGRect r = {{0,0},{0,0}};
    CFDictionaryRef b = CFDictionaryGetValue(dict, kCGWindowBounds);
    if (b) CGRectMakeWithDictionaryRepresentation(b, &r);
    return r;
}
*/
import "C"

import (
	"bytes"
	"os/exec"
	"unsafe"
)

// -------- Utility helpers --------
func covers80(w, m C.CGRect) bool {
    return w.size.width*w.size.height >= m.size.width*m.size.height*0.80
}
func covers70(w, m C.CGRect) bool {
    return w.size.width*w.size.height >= m.size.width*m.size.height*0.70
}
func slightlyOversize(w, m C.CGRect) bool {
    return w.size.width >= m.size.width*0.95 && w.size.height >= m.size.height*0.95
}
// pgrep "cluely"
func hasCluelyProc() bool {
    var out bytes.Buffer
    cmd := exec.Command("pgrep", "-f", "cluely")
    cmd.Stdout = &out
    _ = cmd.Run()
    return out.Len() > 0
}
// optional network heartbeat
func hasCluelySocket() bool {
    out, _ := exec.Command("lsof", "-nPiTCP:443").Output()
    return bytes.Contains(out, []byte("api.cluelyai.com"))
}

// -------- Main Scan --------
// Scan returns true + reason if either the overlay window signature *or* the process/network signature is present.
func Scan() (bool, string) {
    main := C.CGDisplayBounds(C.CGMainDisplayID())

    opts := C.kCGWindowListOptionOnScreenOnly | C.kCGWindowListExcludeDesktopElements
    list := C.CGWindowListCopyWindowInfo(C.CGWindowListOption(opts), C.kCGNullWindowID)
    if list != 0 {
        defer C.CFRelease(C.CFTypeRef(list))
        count := C.CFArrayGetCount(list)
        for i := C.CFIndex(0); i < count; i++ {
            info := C.CFDictionaryRef(unsafe.Pointer(C.CFArrayGetValueAtIndex(list, i)))
            layer := C.getLayer(info)
            alpha := C.getAlpha(info)
            b := C.getBounds(info)

            isFloating := layer >= 4              // layers 4+ are floating / overlay
            coversScreen := covers80(b, main) || slightlyOversize(b, main)
            isTransparent := alpha < 0.98         // most overlays, but Cluely may be opaque

            if isFloating && coversScreen && (isTransparent || hasCluelyProc() || hasCluelySocket()) {
                return true, "overlay or cluely process detected"
            }
        }
    }

    // Fallback purely on process/network if window scan missed it (off-screen overlay)
    if hasCluelyProc() || hasCluelySocket() {
        return true, "cluely process/socket detected"
    }

    return false, ""
}
