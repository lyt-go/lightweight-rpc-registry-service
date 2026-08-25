package service

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNodeCollectionClosesStreamsAndReturnsProducerError(t *testing.T) {
	normalCtx, normalCancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer normalCancel()
	got, err := CollectNodeNames(normalCtx, []string{"edge-a", "edge-b"})
	if err != nil {
		t.Fatalf("normal node collection did not finish: %v", err)
	}
	if want := []string{"edge-a", "edge-b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("normal node collection lost results: got %v want %v", got, want)
	}
	errorCtx, errorCancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer errorCancel()
	_, err = CollectNodeNames(errorCtx, []string{"edge-a", ""})
	if err == nil || !strings.Contains(err.Error(), "node name is empty") {
		t.Fatalf("invalid node should return the producer error without hanging: %v", err)
	}
}
