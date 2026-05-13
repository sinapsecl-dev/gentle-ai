# SDD-CRISP forecasting demo

This demo shows how the experimental `gentle-ai-crisp` fork profile uses **CRISP-DM for data-science governance** and then hands approved engineering work to the existing **SDD implementation chassis**. The important result is not that the toy model wins. The important result is that the workflow can say **no-go**, preserve the evidence, and still produce a controlled SDD handoff when you explicitly choose an engineering exercise.

> Status: experimental fork demo for `sinapsecl-dev/gentle-ai`. This is not an upstream Gentle AI release.

## Quick path

1. Install the experimental CRISP binary into OpenCode.
2. In OpenCode, switch to `crisp-orchestrator-sdd-crisp` with **Shift+Tab**.
3. Ask for a CRISP forecasting MVP over the demo sales dataset.
4. Let CRISP run the phase agents and write artifacts under `crisp/sales-forecast-mvp/`.
5. Review the no-go evaluation: the naive baseline beats the candidate model.
6. If you want to demo the engineering path anyway, explicitly choose the forced SDD handoff option.

## What this proves

| Area | Demo outcome |
|---|---|
| CRISP governance | The workflow evaluates business goal, data, preparation, modeling, evaluation, and deployment readiness before implementation. |
| Honest model decision | The candidate linear regression is rejected because it underperforms the naive baseline. |
| SDD handoff | When forced, CRISP creates SDD spec/design/tasks artifacts instead of implementing inline. |
| Agent separation | CRISP phase agents do data-science planning; SDD agents own implementation work. |
| Safety boundary | A no-go model is not presented as a validated MVP. |

## Install the experimental profile

Build or download the fork-only binary named `gentle-ai-crisp.exe`, then install the OpenCode SDD-CRISP profile:

```powershell
.\gentle-ai-crisp.exe install --agent opencode --component sdd --component skills --sdd-mode multi
```

Expected verification summary:

```text
48 passed, 0 failed
```

After install, OpenCode contains both the original SDD conductor and the CRISP conductor. The installer UI still looks like normal Gentle AI; CRISP is visible in OpenCode through the agent picker.

## Select the CRISP orchestrator

1. Start OpenCode in the project where you want the demo artifacts.
2. Press **Shift+Tab** until `crisp-orchestrator-sdd-crisp` is selected.
3. Keep `gentle-orchestrator` available for normal SDD work; do not replace it.

The CRISP profile is additive. It does not remove the default SDD profile.

## Demo prompt

Use a prompt like this:

```text
Create a CRISP-DM sales forecasting MVP for a tiny daily sales dataset.
Use only date and sales. Exclude promo_flag and weekday.
Compare a naive baseline against a simple candidate model.
Write file-based artifacts so I can show the phase outputs in a demo.
```

For the current demo fixture, the dataset is intentionally tiny. Treat it as a workflow demonstration, not as evidence that the model is production-ready.

## CRISP artifacts to show

During the demo, show the phase outputs as the audit trail:

```text
crisp/sales-forecast-mvp/business-understanding/business-understanding.md
crisp/sales-forecast-mvp/data-understanding/phase-report.md
crisp/sales-forecast-mvp/data-preparation/phase-report.md
crisp/sales-forecast-mvp/modeling/phase-report.md
crisp/sales-forecast-mvp/evaluation/phase-report.md
crisp/sales-forecast-mvp/deployment/phase-report.md
```

These files are the point of CRISP in the demo: they make the data-science decision visible before anyone starts shipping code.

## Expected evaluation result

The real demo run produced a no-go decision:

| Model | MAE |
|---|---:|
| Naive baseline | `4.0` |
| Linear regression candidate | `9.494888834935168` |

Decision: **no-go** for the candidate model. The baseline is better, and the dataset is too small to justify claiming a validated forecasting MVP.

Say this explicitly in the presentation: **the workflow succeeded because it rejected a weak model.** That is the professional behavior. Shipping the candidate as if it were validated would be the failure.

## Forced SDD handoff path

If you want to continue for engineering-demo purposes, choose the explicit forced option. Frame it as:

> We are not approving this model for production. We are forcing an SDD handoff only to demonstrate the controlled implementation path and reproducible reporting.

The CRISP orchestrator should create SDD handoff artifacts like:

```text
openspec/changes/forced-crisp-no-go-handoff/specs/forecasting-handoff/spec.md
openspec/changes/forced-crisp-no-go-handoff/design.md
openspec/changes/forced-crisp-no-go-handoff/tasks.md
```

The handoff must preserve the CRISP constraints:

- use only `date` and `sales`
- exclude `promo_flag` and `weekday`
- use `lag_1`, `lag_2`, and `rolling_mean_3`
- use temporal split only
- reproduce the CRISP metrics
- report `no-go`
- allow baseline fallback/reference behavior
- do not present the candidate as approved for deployment

## Presentation checklist

- [ ] Show that the binary/profile is experimental and fork-only.
- [ ] Show OpenCode agent selection with `crisp-orchestrator-sdd-crisp`.
- [ ] Show the CRISP phase artifacts under `crisp/sales-forecast-mvp/`.
- [ ] Show the evaluation table where baseline beats candidate.
- [ ] State clearly that no-go is a successful governance outcome.
- [ ] If continuing, state that forced SDD handoff is only an engineering exercise.
- [ ] Show the generated SDD handoff artifacts under `openspec/changes/forced-crisp-no-go-handoff/`.
- [ ] Avoid claiming production readiness or model validation.

## Current limitations

| Limitation | Why it matters |
|---|---|
| Tiny dataset | Metrics are useful for demonstrating flow, not model quality. |
| Experimental fork binary | This is not an upstream-supported release yet. |
| OpenCode-focused | The CRISP profile activation path was validated for OpenCode. |
| Installer UI unchanged | CRISP appears in OpenCode after install, not as a separate installer screen. |
| Manual demo flow | Full end-to-end fixture automation is intentionally deferred. |

## Next steps after the demo

- Add CRISP artifact-store selection: Engram, file-based, or hybrid.
- Add automated acceptance coverage for the CRISP no-go handoff path.
- Decide whether the installer TUI should advertise CRISP explicitly.
- If the experiment graduates, add real release packaging for `gentle-ai-crisp` instead of manual fork assets.
