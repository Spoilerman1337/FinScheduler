import {describe, expect, it} from 'vitest';
import {
    buildTagModification,
    createDefaultTagFormData,
    mapTagToFormData,
    normalizeTagFormData,
    runFieldValidators,
    shouldConfirmTagDeactivation,
    tagFormValidators,
    tagNameValidators,
} from './form.ts';

describe('tags form', () => {
    it('creates the default tag form data', () => {
        // Arrange
        const expectedFormData = {
            name: '',
            isActive: true,
        };

        // Act
        const formData = createDefaultTagFormData();

        // Assert
        expect(formData).toEqual(expectedFormData);
    });

    it('maps a tag dto into edit form state', () => {
        // Arrange
        const item = {
            name: 'Transport',
            isActive: false,
        };

        // Act
        const formData = mapTagToFormData(item);

        // Assert
        expect(formData).toEqual({
            name: 'Transport',
            isActive: false,
        });
    });

    it('returns a validation message when the name is blank', () => {
        // Arrange
        const value = '   ';

        // Act
        const validationResult = runFieldValidators(value, tagNameValidators);

        // Assert
        expect(validationResult).toBe('Название обязательно для заполнения');
    });

    it('returns true for a valid name', () => {
        // Arrange
        const value = 'Transport';

        // Act
        const validationResult = tagFormValidators.name(value);

        // Assert
        expect(validationResult).toBe(true);
    });

    it('normalizes tag form data for reuse in the form state', () => {
        // Arrange
        const formData = {
            name: ' Transport ',
            isActive: false,
        };

        // Act
        const normalizedFormData = normalizeTagFormData(formData);

        // Assert
        expect(normalizedFormData).toEqual({
            name: 'Transport',
            isActive: false,
        });
    });

    it('builds a normalized tag modification from valid form data', () => {
        // Arrange
        const formData = {
            name: ' Transport ',
            isActive: false,
        };

        // Act
        const payload = buildTagModification(formData);

        // Assert
        expect(payload).toEqual({
            name: 'Transport',
            isActive: false,
        });
    });

    it('returns true when an active tag is being deactivated in edit mode', () => {
        // Arrange
        const tag = {
            name: 'Transport',
            isActive: true,
        };
        const formData = {
            name: 'Transport',
            isActive: false,
        };

        // Act
        const shouldConfirm = shouldConfirmTagDeactivation('edit', tag, formData);

        // Assert
        expect(shouldConfirm).toBe(true);
    });

    it('returns false when the tag remains active', () => {
        // Arrange
        const tag = {
            name: 'Transport',
            isActive: true,
        };
        const formData = {
            name: 'Transport',
            isActive: true,
        };

        // Act
        const shouldConfirm = shouldConfirmTagDeactivation('edit', tag, formData);

        // Assert
        expect(shouldConfirm).toBe(false);
    });

    it('returns false when creating an inactive tag', () => {
        // Arrange
        const formData = {
            name: 'Transport',
            isActive: false,
        };

        // Act
        const shouldConfirm = shouldConfirmTagDeactivation('create', null, formData);

        // Assert
        expect(shouldConfirm).toBe(false);
    });
});
