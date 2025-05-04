import { createContext, useContext } from "react";
import { PaasibleApi, PaasibleApiCfg } from "./api";

export const PaasibleApiContext = createContext<PaasibleApi>({} as PaasibleApi)

export const PaasibleApiProvider = ({ children, config }: { children: React.ReactNode, config: PaasibleApiCfg }) => {
    const api = PaasibleApi(config)

    return (
        <PaasibleApiContext.Provider value={api}>
            {children}
        </PaasibleApiContext.Provider>
    )
}

export const usePaasibleApi = (): PaasibleApi => {
    const context = useContext(PaasibleApiContext)

    if (!context) {
        throw new Error("usePaasibleApi must be used within a PaasibleApiProvider")
    }

    return context
}
