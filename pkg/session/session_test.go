package session

import (
    "testing"
)

func TestSession_AddFlash_GetFlash(t *testing.T) {
    s := NewSession()

    // GetFlash on empty session returns nil
    flashes := s.GetFlash()
    if flashes != nil {
        t.Errorf("Expected nil flashes on empty session, got %v", flashes)
    }

    // AddFlash then GetFlash
    s.AddFlash("error")
    flashes = s.GetFlash()
    if len(flashes) != 1 || flashes[0] != "error" {
        t.Errorf("Expected one 'error' flash, got %v", flashes)
    }

    // GetFlash consumes the value
    flashes = s.GetFlash()
    if flashes != nil {
        t.Errorf("Expected flash to be consumed after first GetFlash, got %v", flashes)
    }
}

func TestSession_AddFlash_Multiple(t *testing.T) {
    s := NewSession()
    s.AddFlash("first")
    s.AddFlash("second")

    flashes := s.GetFlash()
    if len(flashes) != 2 {
        t.Errorf("Expected 2 flashes, got %d", len(flashes))
    }
    if flashes[0] != "first" || flashes[1] != "second" {
        t.Errorf("Unexpected flash values: %v", flashes)
    }
}

func TestSession_Get_Missing(t *testing.T) {
    s := NewSession()
    v := s.Get("nonexistent")
    if v != nil {
        t.Errorf("Expected nil for missing key, got %v", v)
    }
}

func TestSession_Get_Present(t *testing.T) {
    s := NewSession()
    s.Set("key", "value")
    v := s.Get("key")
    if v != "value" {
        t.Errorf("Expected 'value', got %v", v)
    }
}

func TestSession_Delete(t *testing.T) {
    s := NewSession()
    s.Set("key", "value")
    s.Delete("key")
    v := s.Get("key")
    if v != nil {
        t.Errorf("Expected nil after Delete, got %v", v)
    }

    // Delete of nonexistent key should not panic
    s.Delete("nonexistent")
}
