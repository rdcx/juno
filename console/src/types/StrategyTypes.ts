import type { Field } from "./FieldTypes";
import type { Filter } from "./FilterTypes";
import type { Selector } from "./SelectorTypes";

export interface Strategy {
    id: string;
    name: string;
    selectors: Selector[];
    fields: Field[];
    filters: Filter[];
}

export interface ListResponse {
    status: string;
    message: string;
    strategies: Strategy[];
}

export interface CreateRequest {
    name: string;
}

export interface CreateResponse {
    status: string;
    message: string;
    selector: Strategy;
}

export interface AddSelectorRequest {
    strategy_id: string;
    selector_id: string;
}

export interface AddSelectorResponse {
    status: string;
    message: string;
    selector: Strategy;
}

export interface RemoveSelectorRequest {
    strategy_id: string;
    selector_id: string;
}

export interface RemoveSelectorResponse {
    status: string;
    message: string;
}

export interface AddFieldRequest {
    strategy_id: string;
    field_id: string;
}

export interface AddFieldResponse {
    status: string;
    message: string;
}

export interface RemoveFieldRequest {
    strategy_id: string;
    field_id: string;
}

export interface RemoveFieldResponse {
    status: string;
    message: string;
}

export interface AddFilterRequest {
    strategy_id: string;
    filter_id: string;
}

export interface AddFilterResponse {
    status: string;
    message: string;
}

export interface RemoveFilterRequest {
    strategy_id: string;
    filter_id: string;
}

export interface RemoveFilterResponse {
    status: string;
    message: string;
}