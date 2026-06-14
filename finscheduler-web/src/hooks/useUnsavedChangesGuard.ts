import {useCallback, useEffect, useState} from 'react';
import {
    useBeforeUnload,
    useBlocker,
    useNavigate,
    type NavigateOptions,
    type To,
} from 'react-router-dom';

interface UseUnsavedChangesGuardOptions {
    isDirty: boolean;
    isDisabled?: boolean;
}

interface PendingNavigation {
    options?: NavigateOptions;
    to: To;
}

export function useUnsavedChangesGuard({
    isDirty,
    isDisabled = false,
}: UseUnsavedChangesGuardOptions) {
    const navigate = useNavigate();
    const [pendingNavigation, setPendingNavigation] = useState<PendingNavigation | null>(null);

    const shouldBlock = !isDisabled && isDirty;

    const blocker = useBlocker(shouldBlock);

    useEffect(() => {
        if (pendingNavigation === null || isDirty) {
            return;
        }

        const {to, options} = pendingNavigation;

        setPendingNavigation(null);
        navigate(to, options);
    }, [isDirty, navigate, pendingNavigation]);

    useBeforeUnload(
        useCallback(
            (event) => {
                if (!shouldBlock) {
                    return;
                }

                event.preventDefault();
            },
            [shouldBlock],
        ),
        {capture: true},
    );

    const scheduleNavigation = useCallback((to: To, options?: NavigateOptions) => {
        setPendingNavigation({to, options});
    }, []);

    const leavePage = useCallback(() => {
        if (blocker.state !== 'blocked') {
            return;
        }

        blocker.proceed();
    }, [blocker]);

    const stayOnPage = useCallback(() => {
        if (blocker.state !== 'blocked') {
            return;
        }

        blocker.reset();
    }, [blocker]);

    return {
        isDialogOpen: blocker.state === 'blocked',
        leavePage,
        scheduleNavigation,
        stayOnPage,
    };
}
