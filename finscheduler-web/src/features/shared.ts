import type {Lookup} from '../api/types.ts';
import type {SelectOption} from '../components/formFields/types.ts';

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
