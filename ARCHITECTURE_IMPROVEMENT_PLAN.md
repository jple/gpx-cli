# Architecture Improvement Plan for GPX-CLI

This document outlines a plan to implement **dependency injection** and a **builder pattern** for the `GeoStatistic` struct in the GPX-CLI codebase. These improvements aim to enhance modularity, testability, and maintainability while preserving the existing functionality.

---

## 1. Dependency Injection

### **Goal**
Decouple the `stat.go` logic from the `geostat` package by introducing an interface for exporting statistics. This will allow for:
- Easier mocking in tests.
- Swapping export implementations without modifying `stat.go`.
- Better adherence to the Dependency Inversion Principle (DIP).

### **Steps**

#### **Step 1: Define an Interface for Exporting Statistics**
Create a new interface `StatExporter` in `internal/gpx/stat.go`:

```go
type StatExporter interface {
    WriteHeader(w io.Writer, kind string, opt geostat.Option) error
    WriteValues(w io.Writer, i *int, kind string, opts ...geostat.Option) error
    WriteFooter(w io.Writer, kind string) error
}
```

#### **Step 2: Update `ExportStatsTrks` to Accept an Exporter**
Modify the `ExportStatsTrks` function to accept a `StatExporter` as a parameter:

```go
func ExportStatsTrks(w io.Writer, trks []Trk, flatspeed float64, detail bool, kind string, exporter StatExporter) {
    // Use exporter instead of direct calls to geostat
    // Example:
    // exporter.WriteHeader(w, kind, exportOpt)
}
```

#### **Step 3: Create a Concrete Implementation**
Implement the `StatExporter` interface in the `geostat` package:

```go
// In internal/gpx/geostat/export.go
type GeoStatExporter struct{}

func (e *GeoStatExporter) WriteHeader(w io.Writer, kind string, opt Option) error {
    // Existing logic
}

func (e *GeoStatExporter) WriteValues(w io.Writer, i *int, kind string, opts ...Option) error {
    // Existing logic
}

func (e *GeoStatExporter) WriteFooter(w io.Writer, kind string) error {
    // Existing logic
}
```

#### **Step 4: Update `info.go` to Inject the Exporter**
Modify the `Run` function in `cmd/info.go` to inject the `GeoStatExporter`:

```go
func (cmd *cobra.Command, args []string) {
    // ... existing code ...
    exporter := &geostat.GeoStatExporter{}
    gpx.ExportStatsTrks(os.Stdout, g.Trks, viper.GetFloat64("speed"), detail.Value(), viper.GetString("kind"), exporter)
}
```

#### **Step 5: Update Tests**
- Mock the `StatExporter` interface in tests to isolate `stat.go` logic.
- Ensure all existing tests pass with the new dependency injection.

---

## 2. Builder Pattern for `GeoStatistic`

### **Goal**
Encapsulate the construction of the `GeoStatistic` struct to:
- Reduce boilerplate code for creating and populating the struct.
- Ensure consistency in how `GeoStatistic` objects are created.
- Make the code more readable and maintainable.

### **Steps**

#### **Step 1: Create a Builder Struct**
Add a `GeoStatisticBuilder` struct and methods to `internal/gpx/geostat/geo_statistic.go`:

```go
type GeoStatisticBuilder struct {
    stat GeoStatistic
}

func NewGeoStatisticBuilder() *GeoStatisticBuilder {
    return &GeoStatisticBuilder{
        stat: GeoStatistic{
            From: "start",
            To:   "end",
        },
    }
}
```

#### **Step 2: Add Builder Methods**
Add methods to set fields on the `GeoStatistic` struct:

```go
func (b *GeoStatisticBuilder) WithName(name string) *GeoStatisticBuilder {
    b.stat.Name = name
    return b
}

func (b *GeoStatisticBuilder) WithFrom(from string) *GeoStatisticBuilder {
    b.stat.From = from
    return b
}

func (b *GeoStatisticBuilder) WithTo(to string) *GeoStatisticBuilder {
    b.stat.To = to
    return b
}

func (b *GeoStatisticBuilder) WithDistance(distance float64) *GeoStatisticBuilder {
    b.stat.Distance = distance
    return b
}

func (b *GeoStatisticBuilder) WithAscentDescent(ascent, descent float64) *GeoStatisticBuilder {
    b.stat.TotalAscent = ascent
    b.stat.TotalDescent = descent
    return b
}

func (b *GeoStatisticBuilder) WithDuration(hour, min int8) *GeoStatisticBuilder {
    b.stat.DurationHour = hour
    b.stat.DurationMin = min
    return b
}

func (b *GeoStatisticBuilder) WithDistanceEffort(effort float64) *GeoStatisticBuilder {
    b.stat.DistanceEffort = effort
    return b
}

func (b *GeoStatisticBuilder) WithN(n int) *GeoStatisticBuilder {
    b.stat.N = n
    return b
}

func (b *GeoStatisticBuilder) Build() GeoStatistic {
    return b.stat
}
```

#### **Step 3: Update `stat.go` to Use the Builder**
Replace direct assignments to `GeoStatistic` with the builder pattern in `stat.go`:

```go
func (trkpts Trkpts) Stats(flatSpeed float64) geostat.GeoStatistic {
    builder := geostat.NewGeoStatisticBuilder()
    builder.WithN(len(trkpts))
    builder.WithDistance(trkpts.TotalDistance())
    builder.WithAscentDescent(
        trkpts.TotalAscent(RollingWindowSize),
        trkpts.TotalDescent(RollingWindowSize),
    )
    
    if trkpts[0].Name != nil {
        builder.WithName(*trkpts[0].Name)
        builder.WithFrom(*trkpts[0].Name)
    }
    
    if trkpts[len(trkpts)-1].Name != nil {
        builder.WithTo(*trkpts[len(trkpts)-1].Name)
    }
    
    stat := builder.Build()
    
    // Set calculation value
    stat.DistanceEffort = geo.CalcDistanceEffort(
        stat.Distance,
        stat.TotalAscent,
        stat.TotalDescent)
    _, stat.DurationHour, stat.DurationMin = geo.CalcDuration(stat.DistanceEffort, flatSpeed)
    
    return stat
}
```

#### **Step 4: Update `StatsPerSection` and Other Functions**
Apply the same builder pattern to other functions like `StatsPerSection`:

```go
func (trk Trk) StatsPerSection(flatSpeed float64) []geostat.GeoStatistic {
    stats := []geostat.GeoStatistic{}
    TrkSections := trk.SplitByTrkptName()
    for i, trkpts := range TrkSections {
        if len(trkpts) == 0 {
            continue
        }
        
        stat := trkpts.Stats(flatSpeed)
        builder := geostat.NewGeoStatisticBuilder().
            WithN(stat.N).
            WithDistance(stat.Distance).
            WithAscentDescent(stat.TotalAscent, stat.TotalDescent).
            WithDistanceEffort(stat.DistanceEffort).
            WithDuration(stat.DurationHour, stat.DurationMin)
        
        if i < len(TrkSections)-1 {
            builder.WithN(stat.N - 1)
        }
        
        // Set From and To
        if i == 0 && trk.Name != "" {
            builder.WithFrom(trk.Name)
        }
        if trkpts[0].Name != nil {
            builder.WithFrom(*trkpts[0].Name)
        }
        if i < len(TrkSections)-1 && len(trkpts) > 0 {
            builder.WithTo(*trkpts[len(trkpts)-1].Name)
        }
        
        stats = append(stats, builder.Build())
    }
    return stats
}
```

#### **Step 5: Update Tests**
- Ensure all tests are updated to use the builder pattern where applicable.
- Verify that the builder produces the same results as the direct construction.

---

## 3. Testing Plan

### **Dependency Injection Tests**
1. **Mock Exporter**:
   - Create a mock implementation of `StatExporter` for testing.
   - Verify that `ExportStatsTrks` correctly calls the exporter methods.

2. **Integration Test**:
   - Test the full flow from `info.go` to `stat.go` with the injected exporter.
   - Ensure output formats (`string`, `csv`, `html`) are generated correctly.

### **Builder Pattern Tests**
1. **Unit Tests**:
   - Test each builder method to ensure it sets the correct field.
   - Verify that `Build()` returns a correctly populated `GeoStatistic`.

2. **Integration Tests**:
   - Test `Stats` and `StatsPerSection` with the builder.
   - Compare results with the previous implementation to ensure consistency.

---

## 4. Migration Plan

### **Phase 1: Dependency Injection**
1. Define the `StatExporter` interface.
2. Update `ExportStatsTrks` to accept the interface.
3. Implement `GeoStatExporter` in the `geostat` package.
4. Update `info.go` to inject the exporter.
5. Test thoroughly.

### **Phase 2: Builder Pattern**
1. Implement the `GeoStatisticBuilder`.
2. Update `Stats` and `StatsPerSection` to use the builder.
3. Test thoroughly.

### **Phase 3: Cleanup**
1. Remove redundant code or comments.
2. Update documentation to reflect the new patterns.
3. Ensure all tests pass.

---

## 5. Risks and Mitigations

| Risk                                  | Mitigation                                                                 |
|---------------------------------------|-----------------------------------------------------------------------------|
| Breaking existing functionality        | Write comprehensive tests before and after changes.                       |
| Increased complexity                   | Keep the builder and interface simple; document thoroughly.                |
| Performance overhead                   | Profile before/after to ensure no significant performance degradation.     |
| Difficulty in debugging                | Add detailed logs or comments for complex builder chains.                  |

---

## 6. Timeline

| Task                                  | Estimated Time |
|---------------------------------------|-----------------|
| Define `StatExporter` interface        | 1 hour          |
| Update `ExportStatsTrks`              | 1 hour          |
| Implement `GeoStatExporter`           | 2 hours         |
| Update `info.go`                       | 1 hour          |
| Test dependency injection             | 2 hours         |
| Implement `GeoStatisticBuilder`       | 2 hours         |
| Update `Stats` and `StatsPerSection`   | 3 hours         |
| Test builder pattern                  | 2 hours         |
| Cleanup and documentation             | 1 hour          |

**Total Estimated Time**: ~15 hours

---

## 7. Open Questions

1. Should the builder pattern be extended to other structs (e.g., `Trk`, `Trkpts`)?
2. Are there other dependencies in `stat.go` that could benefit from injection?
3. Should the `StatExporter` interface include additional methods for future formats (e.g., JSON)?

---

## 8. Next Steps

1. Review this plan with the team and gather feedback.
2. Prioritize tasks based on impact and feasibility.
3. Begin implementation with dependency injection, followed by the builder pattern.
4. Run tests after each phase to ensure stability.

---

**Note**: This plan assumes the current architecture and requirements remain unchanged. Adjust as needed based on feedback or new insights.