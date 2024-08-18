/**
* This file was generated by mirror, do not edit it manually as it will be overwritten.
*
* You can find the docs and source code for mirror here: https://github.com/aosasona/mirror
*/


type Flattened_Language = string;

type Flattened_Time = number;

type Flattened_Tags = Record<string, string>;

type Flattened_Person = {
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
    created_at: Time;
    updated_at: number | null;
    deleted_at: Time;
    is_active: boolean;
};

type Flattened_CreateUserFunc = (arg0: Person) => string;