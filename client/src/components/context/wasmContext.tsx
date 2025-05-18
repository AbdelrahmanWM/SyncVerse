import  {createContext, useContext} from 'react';
import {IWasmContext} from "../../interfaces/IWasmContext";

export const WasmContext = createContext<IWasmContext>({});

export const useWasm = ()=>{
    const context=  useContext(WasmContext);
    if(!context){
        throw new Error('useWasm must be used within a WasmProvider');
    }
    return context;
};