package astilectron

import "github.com/asticode/go-astitools/context"

// Tray event names
const (
	EventNameTrayCmdCreate      = "tray.cmd.create"
	EventNameTrayCmdDestroy     = "tray.cmd.destroy"
	EventNameTrayEventCreated   = "tray.event.created"
	EventNameTrayEventDestroyed = "tray.event.destroyed"
)

// Tray represents a tray
type Tray struct {
	*object
	m *Menu
	o *TrayOptions
}

// TrayOptions represents tray options
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use PtrBool, PtrInt or PtrStr
// to fill the struct
// https://github.com/electron/electron/blob/v1.6.5/docs/api/tray.md
type TrayOptions struct {
	Image   *string `json:"image,omitempty"`
	Tooltip *string `json:"tooltip,omitempty"`
}

// newTray creates a new tray
func newTray(o *TrayOptions, c *asticontext.Canceller, d *Dispatcher, i *identifier, wrt *writer) (t *Tray) {
	// Init
	t = &Tray{
		o:      o,
		object: newObject(nil, c, d, i, wrt),
	}

	// Make sure the tray's context is cancelled once the destroyed event is received
	t.On(EventNameTrayEventDestroyed, func(e Event) (deleteListener bool) {
		t.cancel()
		return true
	})
	return
}

// Create creates the tray
func (t *Tray) Create() (err error) {
	if err = t.isActionable(); err != nil {
		return
	}
	var e = Event{Name: EventNameTrayCmdCreate, TargetID: t.id, TrayOptions: t.o}
	if t.m != nil {
		e.MenuID = t.m.id
	}
	_, err = synchronousEvent(t.c, t, t.w, e, EventNameTrayEventCreated)
	return
}

// Destroy destroys the tray
func (t *Tray) Destroy() (err error) {
	if err = t.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(t.c, t, t.w, Event{Name: EventNameTrayCmdDestroy, TargetID: t.id}, EventNameTrayEventDestroyed)
	return
}

// NewMenu creates a new tray menu
func (t *Tray) NewMenu(i []*MenuItemOptions) *Menu {
	t.m = newMenu(t.ctx, t, i, t.c, t.d, t.i, t.w)
	return t.m
}
