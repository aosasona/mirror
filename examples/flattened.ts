/**
 * This file was generated by mirror, do not edit it manually as it will be overwritten.
 *
 * You can find the docs and source code for mirror here: https://github.com/aosasona/mirror
 */

export type Flattened_Language = string;

export type Flattened_Tags = Record<string, string>;

export type Flattened_Person = {
    first_name: string;
    last_name: string;
    age: number;
    line_1: string | null;
    line_2: string | null;
    street: string;
    city: string;
    state: string;
    postal_code: string;
    country: string;
    languages: Array<string>;
    grades?: Record<string, number>;
    tags: Record<string, string>;
    created_at: string;
    updated_at: number | null;
    deleted_at: string | null;
    is_active: boolean;
};

export type Flattened_Collection = {
    items: Array<string>;
    desc: string;
};

export type Flattened_CreateUserFunc = (arg0: Person) => string;