# Bridge Removal Criteria

The `seer-q agent install/uninstall/doctor` compatibility bridge can be removed when ALL of the following are verified.

## Prerequisites

1. [ ] Private skills repository (`SparkssL/seer-skills`) exists and is published
2. [ ] `npx skills add SparkssL/seer-skills --all -y` works with GitHub auth
   - Test with `GITHUB_TOKEN`
   - Test with SSH key auth
   - Verify skills land in `.claude/skills/*/SKILL.md`
3. [ ] At least one non-bridge install path verified end-to-end for Claude Code
4. [ ] Published skills repo README has working install instructions
5. [ ] Skills repo includes `commands/claude/seer.md` for optional `/seer` command

## Cleanup checklist (execute after prerequisites pass)

- [ ] Delete `internal/cmd/agent/` (6 files)
- [ ] Delete `agent/` (embed.go + cmd/seer.md)
- [ ] Delete `skills/embed.go`
- [ ] Remove agent entry from `internal/registry/registry.go`
- [ ] Remove agent import from registry.go
- [ ] Update `consistency_test.go` — read skills from disk, remove cmd test
- [ ] Remove agent checks from `test/smoke-test.sh`
- [ ] Remove agent section from `docs/cli-reference.md`
- [ ] Remove bridge checks from `docs/release-gate.md`
- [ ] Update `CHANGELOG.md`
- [ ] Delete this file
