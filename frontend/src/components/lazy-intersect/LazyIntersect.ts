export interface IntersectEvent extends CustomEvent {
    detail: IntersectionObserverEntry;
}

function createIntersectionObserver() {
    const intersectionObserver = new IntersectionObserver(
        (entries) =>
            entries.forEach((entry) =>
                entry.target.dispatchEvent(
                    new CustomEvent("intersect", { detail: entry })
                )
            ),
        { root: null, rootMargin: "0px", threshold: 0 }
    );
    return intersectionObserver;
}

function createMutationObserver(intersectionObserver: IntersectionObserver) {
    const mutationObserver = new MutationObserver((mutations) =>
        mutations.forEach((m) => {
            m.addedNodes.forEach((node) => {
                if (
                    node instanceof HTMLElement &&
                    node.dataset.intersect != null &&
                    node.dataset.intersectInitialized == null
                ) {
                    intersectionObserver.observe(node);
                    node.dataset.intersectInitialized = "true";
                }
            });
            m.removedNodes.forEach((node) => {
                if (node instanceof HTMLElement) {
                    intersectionObserver.unobserve(node);
                }
            });
        })
    );

    return mutationObserver;
}

/**
 * Creates an IntersectionObserver and MutationObserver, which together make lazy loading possible for any HTML element with the `data-intersect` prop and `on:intersect` event.
 * MutationObserver checks for any new data-intersect nodes and adds them to the IntersectionObserver; removed nodes are unobserved.
 * IntersectionObserver dispatches the intersect event and passes the IntersectionObserverEntry event from IntersectionObserver.
 * on:intersect can be used on the node to call appropriate UI functions.
 * A new 'HTMLAttributes<T>' interface needs to be created on the 'svelteHTML' global scope in app.d.ts; it must include an event type for 'on:intersect'.
 * @returns
 */
export function createLazyIntersect(): {
    intersectionObserver: IntersectionObserver;
    mutationObserver: MutationObserver;
} {
    const intersectionObserver = createIntersectionObserver();
    const mutationObserver = createMutationObserver(intersectionObserver);

    [...document.querySelectorAll("[data-intersect]")].forEach((node: any) => {
        intersectionObserver.observe(node);
        node.dataset.intersectInitialized = "true";
    });

    mutationObserver.observe(document.body, {
        childList: true,
        subtree: true,
    });

    return { intersectionObserver, mutationObserver };
}
