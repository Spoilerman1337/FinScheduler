import {describe, expect, it} from 'vitest';
import {
    buildTagModification,
    createDefaultTagFormData,
    mapTagToFormData,
    shouldConfirmTagDeactivation,
    validateTagFormData,
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

    it('returns an error when the name is blank', () => {
        // Arrange
        const formData = {
            ...createDefaultTagFormData(),
            name: '   ',
        };

        // Act
        const error = validateTagFormData(formData);

        // Assert
        expect(error).toBe('Название обязательно для заполнения');
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
