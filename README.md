# go-integer-range

In golang, converts classic for loops (`for i := 0; i<10; i++`) into integer rangers (`for i := range 10`) introduced in go 1.22.

# Installation

```
go install github.com/lukemassa/go-integer-range@v0.2.2
go-integer-range -h
```

If `GOBIN` is not in your path, try `~/go/bin/go-integer-range -h`

# Usage

## Dryrun update of a whole directory
```
documentation-maintainer % ~/go/bin/go-integer-range --dryrun pkg/docs
2024/03/01 12:23:54 Would have updated pkg/docs/collector.go, skipping for dryrun
2024/03/01 12:23:54 Would have updated pkg/docs/collector_test.go, skipping for dryrun
2024/03/01 12:23:54 Would have updated pkg/docs/config.go, skipping for dryrun
2024/03/01 12:23:54 No updates needed for pkg/docs/config_test.go
2024/03/01 12:23:54 Would have updated pkg/docs/confluence_space.go, skipping for dryrun
2024/03/01 12:23:54 No updates needed for pkg/docs/gitlab_group.go
2024/03/01 12:23:54 No updates needed for pkg/docs/gitlab_group_test.go
2024/03/01 12:23:54 Would have updated pkg/docs/jira.go, skipping for dryrun
2024/03/01 12:23:54 No updates needed for pkg/docs/jira_test.go
2024/03/01 12:23:54 No updates needed for pkg/docs/maintainer.go
2024/03/01 12:23:54 Would have updated pkg/docs/maintainer_test.go, skipping for dryrun
2024/03/01 12:23:54 No updates needed for pkg/docs/repo.go
2024/03/01 12:23:54 No updates needed for pkg/docs/repo_test.go
```

## Dryrun update of a single file
```
documentation-maintainer % ~/go/bin/go-integer-range --dryrun pkg/docs/collector.go 
2024/03/01 12:24:01 Would have updated pkg/docs/collector.go, skipping for dryrun
```

## Actual update of a single file

```
documentation-maintainer % ~/go/bin/go-integer-range pkg/docs/collector.go         
2024/03/01 12:24:05 Updating pkg/docs/collector.go

documentation-maintainer % git diff pkg/docs/collector.go        
diff --git a/pkg/docs/collector.go b/pkg/docs/collector.go
index 3e8c886..7015343 100644
--- a/pkg/docs/collector.go
+++ b/pkg/docs/collector.go
@@ -65,22 +65,20 @@ func (d DocumentCollection) Get() []Document {
 // Split up the docs into num equal piles
 func (d DocumentCollection) Divvy(num int) [][]Document {
        ret := make([][]Document, num)
-
-       for i := 0; i < num; i++ {
+       for i := range num {
                ret[i] = make([]Document, 0)
        }
```

Note that the script is idempotent:
```
documentation-maintainer % ~/go/bin/go-integer-range pkg/docs/collector.go
2024/03/01 12:24:07 No updates needed for pkg/docs/collector.go
```