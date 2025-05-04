import { expandedFetch } from "../http";
import { PaasiblePB } from "./pocketbase";
import { uuidv7 } from "uuidv7";

export type PaasibleApiCfg = {
  host: string;
  pb: PaasiblePB;
};

export type GetLocalMachineIdResponseData = {
  value: string;
};

export type CheckLocalMachineId =
  | {
      status: "no_name";
      value: "";
    }
  | {
      status: "not_found";
      value: string;
    }
  | {
      status: "found";
      value: string;
    };

export type GetRepositoriesResponseData = {
  repositories: string[];
};

export type PaasibleApi = ReturnType<typeof PaasibleApi>;
export const PaasibleApi = (api: PaasibleApiCfg) => {
  const { pb } = api;

  return {
    pb,
    Queries: {
      getLocalMachineId: () =>
        expandedFetch<GetLocalMachineIdResponseData>(
          fetch(`${api.host}/get-local-machine-id`)
        ),
      checkLocalMachineId: () =>
        expandedFetch<CheckLocalMachineId>(
          fetch(`${api.host}/check-local-machine-id`)
        ),
      getRepositories: () => pb.collection("repository").getFullList(),
    },

    Mutations: {
      signIn: async (values: { email: string; password: string }) => {
        return pb
          .collection("users")
          .authWithPassword(values.email, values.password);
      },
      createCurrentMachine: (machineId: string) => {
        return pb.collection("machine").create({
          name: machineId,
          current: true,
          id: uuidv7(),
        });
      },
    },
  };
};
