import PocketBase, { RecordService } from "pocketbase";
import { createContext, useContext } from "react";

export type Base = {
  id: string;
};

export type BaseWithTimestamp = Base & {
  created: string;
  updated: string;
};

export type Repository = BaseWithTimestamp & {
  remote: string;
  top_level: string;
};

export type Machine = BaseWithTimestamp & {
  name: string;
};

export type Playbook = BaseWithTimestamp & {
  name: string;
  path: string;
  content: string;
  repository_id: string;
};

export type PlaybookRunResult = BaseWithTimestamp & {
  command: string;
  stdout: string;
  stderr: string;
  error: string;
  pwd: string;
  repository_id: string;
  repository_branch: string;
  playbook_id: string;
  machine_id: string;
};

export type PaasiblePB = PocketBase & {
  collection(idOrName: string): RecordService; // default fallback for any other collection
  collection(idOrName: "repository"): RecordService<Repository>;
  collection(idOrName: "machine"): RecordService<Machine>;
  collection(idOrName: "playbook"): RecordService<Playbook>;
  collection(idOrName: "playbook_run_result"): RecordService<PlaybookRunResult>;
};

export const PocketBaseContext = createContext<PocketBase>({} as PocketBase);

export const PocketBaseProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  // TODO: Make this configurable
  const pb = new PocketBase("http://127.0.0.1:8090") as PaasiblePB;

  return (
    <PocketBaseContext.Provider value={pb}>
      {children}
    </PocketBaseContext.Provider>
  );
};

export const usePocketBase = (): PocketBase => {
  const context = useContext(PocketBaseContext);
  if (!context) {
    throw new Error("usePocketBase must be used within a PocketBaseProvider");
  }
  return context;
};
