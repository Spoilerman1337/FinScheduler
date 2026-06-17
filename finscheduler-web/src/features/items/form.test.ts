import {describe, expect, it} from 'vitest';
import {
    buildItemModification,
    createDefaultItemFormData,
    itemCashbackValidators,
    itemCategoryValidators,
    itemFormValidators,
    itemNameValidators,
    itemPriceValidators,
    mapItemToFormData,
    normalizeItemFormData,
} from './form.ts';
import {runFieldValidators} from '../shared.ts';

describe('items form', () => {
    it('creates the default item form data', () => {
        // Arrange
        const expectedFormData = {
            name: '',
            description: '',
            price: '',
            cashback: '',
            isActive: true,
            category: '',
            tagIds: [],
        };

        // Act
        const formData = createDefaultItemFormData();

        // Assert
        expect(formData).toEqual(expectedFormData);
    });

    it('maps an item dto into edit form state', () => {
        // Arrange
        const item = {
            name: 'Coffee',
            description: 'Morning coffee',
            price: 199.5,
            cashback: 5,
            isActive: false,
            category: 'FoodDrinks',
            tags: [{value: 'tag-1', label: 'Food'}, {value: 'tag-2'}, {value: ''}],
            priceHistory: [],
            priceForecast: [],
        };

        // Act
        const formData = mapItemToFormData(item);

        // Assert
        expect(formData).toEqual({
            name: 'Coffee',
            description: 'Morning coffee',
            price: '199.5',
            cashback: '5',
            isActive: false,
            category: 'FoodDrinks',
            tagIds: ['tag-1', 'tag-2'],
        });
    });

    it('returns a validation message when the name is blank', () => {
        // Arrange
        const value = '   ';

        // Act
        const validationResult = runFieldValidators(value, itemNameValidators);

        // Assert
        expect(validationResult).toBe('Название обязательно для заполнения');
    });

    it('returns a validation message when the price is not numeric', () => {
        // Arrange
        const value = 'abc';

        // Act
        const validationResult = runFieldValidators(value, itemPriceValidators);

        // Assert
        expect(validationResult).toBe('Цена должна быть числом');
    });

    it('returns a validation message when cashback is out of range', () => {
        // Arrange
        const value = '101';

        // Act
        const validationResult = runFieldValidators(value, itemCashbackValidators);

        // Assert
        expect(validationResult).toBe('Кэшбэк должен быть числом от 0 до 100');
    });

    it('returns a validation message when the category is missing or unknown', () => {
        // Arrange
        const value = 'UnknownCategory';

        // Act
        const validationResult = runFieldValidators(value, itemCategoryValidators);

        // Assert
        expect(validationResult).toBe('Выберите категорию');
    });

    it('returns true for valid field values', () => {
        // Arrange
        const validValues = {
            name: 'Coffee',
            price: '',
            cashback: '10',
            category: 'FoodDrinks',
        };

        // Act
        const validationResults = {
            name: itemFormValidators.name(validValues.name),
            price: itemFormValidators.price(validValues.price),
            cashback: itemFormValidators.cashback(validValues.cashback),
            category: itemFormValidators.category(validValues.category),
        };

        // Assert
        expect(validationResults).toEqual({
            name: true,
            price: true,
            cashback: true,
            category: true,
        });
    });

    it('normalizes item form data for reuse in the form state', () => {
        // Arrange
        const formData = {
            name: ' Coffee ',
            description: ' Morning coffee ',
            price: ' 199.5 ',
            cashback: ' 3 ',
            isActive: true,
            category: 'FoodDrinks',
            tagIds: ['tag-1', 'tag-2'],
        };

        // Act
        const normalizedFormData = normalizeItemFormData(formData);

        // Assert
        expect(normalizedFormData).toEqual({
            name: 'Coffee',
            description: 'Morning coffee',
            price: '199.5',
            cashback: '3',
            isActive: true,
            category: 'FoodDrinks',
            tagIds: ['tag-1', 'tag-2'],
        });
    });

    it('builds a normalized item payload from valid form data', () => {
        // Arrange
        const formData = {
            name: ' Coffee ',
            description: ' Morning coffee ',
            price: '',
            cashback: '3',
            isActive: true,
            category: 'FoodDrinks',
            tagIds: ['tag-1', 'tag-2'],
        };

        // Act
        const payload = buildItemModification(formData);

        // Assert
        expect(payload).toEqual({
            name: 'Coffee',
            description: 'Morning coffee',
            price: 0,
            cashback: 3,
            isActive: true,
            category: 'FoodDrinks',
            tagIds: ['tag-1', 'tag-2'],
        });
    });
});
