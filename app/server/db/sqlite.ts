import { Database } from "bun:sqlite";
import { dirname } from "node:path";
import { mkdirSync } from "node:fs";
import { schemaStatements } from "./schema";

let database: Database | null = null;

export function getDatabase(): Database {
  if (database) {
    return database;
  }

  const path = process.env.SQLITE_PATH || `${import.meta.dir}/../../../data/pos.sqlite`;
  mkdirSync(dirname(path), { recursive: true });
  database = new Database(path, { create: true });
  database.exec("PRAGMA journal_mode = WAL;");
  database.exec("PRAGMA synchronous = NORMAL;");
  database.exec("PRAGMA busy_timeout = 5000;");

  for (const statement of schemaStatements) {
    database.exec(statement);
  }

  return database;
}

export function closeDatabase(): void {
  if (!database) {
    return;
  }

  database.close(false);
  database = null;
}
