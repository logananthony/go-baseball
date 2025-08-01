# Modeling Overview

This document provides a high level look at how the simulation engine models a baseball game. It is intended for users who want a general idea of the workflow and the types of data that drive the outcomes.

## Data Preparation

1. **Game context** – `PrepareSimData` collects the game information for the desired `gamePk`, including lineups and starting pitchers. It also loads park factors, bullpen orders and substitution probabilities.
2. **Player metrics** – For every batter and pitcher in the lineups, a variety of historical metrics are fetched from the database using the fetcher package. These include:
   - Batter swing and contact rates
   - Batter hit type probabilities (single, double, etc.)
   - Pitcher pitch usage by count
   - Pitcher location and velocity distributions
   - League‐level swing/contact distributions as fallbacks
3. **Assembling `SimData`** – The collected records are bundled into a `SimData` struct that is passed to the game simulator. `SimData` holds all player info, league averages and park factor adjustments needed during simulation.

## Game Simulation

`SimulateGame` runs a full nine inning game and uses `SimulatePlateAppearance` to handle each at‑bat. The function keeps track of the lineup order, base state, score and pitching substitutions.

### Plate Appearances

`SimulatePlateAppearance` produces a sequence of pitch outcomes until the plate appearance resolves. It follows these general steps on every pitch:

1. **Pitch Type Selection** – `SimulatePitchType` samples a pitch type based on the pitcher’s historical pitch mix for the current ball/strike count. If no direct data is available, nearby counts or overall frequencies are used.
2. **Location & Velocity** – `SimulatePitchLocationVelo` draws a location and velocity from a multivariate normal distribution. Player level parameters are used when there is sufficient data; otherwise league averages act as a fallback.
3. **Swing Decision** – Using batter swing percentages (`SimulateSwingDecision`) the model decides whether the batter offers at the pitch.
4. **Contact Result** – If the batter swings, `SimulateContactPercentage` determines whether the result is a whiff, a foul or a ball in play.
5. **Ball In Play** – When the ball is put in play, `SimulateBatterHitType` chooses the specific event (single, double, etc.). Probabilities are adjusted by park factors for the home team’s stadium.

These events are recorded in a `PlateAppearanceResult`, and `ProcessPlateAppearance` updates the base state and score accordingly.

### Pitching Changes

During the game the simulator monitors pitch counts and score differentials. When thresholds are met, bullpen role probabilities are consulted to select a relief pitcher from the available bullpen lineup.

## Aggregation and Output

Each completed game yields a `GameResult` with detailed pitch and event data. Higher level endpoints in the API run many simulations in parallel and aggregate the results into probability summaries (run totals, strikeout props, etc.) for users.

---

This is only a summary of the main modeling flow. The source under `pkg/sim` contains the exact implementation details.
