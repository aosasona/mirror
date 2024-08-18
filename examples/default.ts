/**
* This file was generated by mirror, do not edit it manually as it will be overwritten.
*
* You can find the docs and source code for mirror here: https://github.com/aosasona/mirror
*/


type Language = string;

type Address = {
	line_1: string | null;
	line_2: string | null;
	street: string;
	city: string;
	state: string;
	postal_code: string;
	country: string;
};

type Tags = Record<string, string>;

type Person = {
	first_name: string;
	last_name: string;
	age: number;
	address: Address;
	languages: Array<string>;
	grades?: Record<string, number>;
	tags: Record<string, string>;
	created_at: string;
	updated_at: number | null;
	deleted_at: string | null;
	is_active: boolean;
};

type CreateUserFunc = (arg0: Person) => string;