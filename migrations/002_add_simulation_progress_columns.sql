ALTER TABLE simulation_jobs
    ADD COLUMN current_simulation INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN total_simulations   INTEGER NOT NULL DEFAULT 0;
