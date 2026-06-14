import type {Lookup} from '../api/types.ts';
import type {SelectOption} from '../components/formFields/types.ts';

export interface FieldValidator<TValue> {
    validate: (value: TValue) => boolean;
    errorMessage: string;
}

export function runFieldValidators<TValue>(
    value: TValue,
    validators: FieldValidator<TValue>[],
): true | string {
    for (const validator of validators) {
        if (!validator.validate(value)) {
            return validator.errorMessage;
        }
    }

    return true;
}

export function mapLookupsToSelectOptions(lookups?: Lookup[] | null): SelectOption[] {
    if (!Array.isArray(lookups)) {
        return [];
    }

    return lookups.reduce<SelectOption[]>((result, lookup) => {
        if (!lookup?.value) {
            return result;
        }

        result.push({
            label: lookup.label ?? lookup.value,
            value: lookup.value,
        });

        return result;
    }, []);
}
