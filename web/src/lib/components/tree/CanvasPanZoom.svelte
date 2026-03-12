<script lang="ts">
    import { onMount } from 'svelte';
    import * as d3 from 'd3';

    let { children, width = '100%', height = '100%' } = $props();

    let containerElement: HTMLDivElement;
    let svgElement: SVGSVGElement;
    let gElement: SVGGElement;

    // Zoom behavior state
    let zoom: d3.ZoomBehavior<SVGSVGElement, unknown>;

    onMount(() => {
        zoom = d3.zoom<SVGSVGElement, unknown>()
            .scaleExtent([0.1, 3])
            .on('zoom', (event) => {
                d3.select(gElement).attr('transform', event.transform);
            });

        d3.select(svgElement).call(zoom);

        // Initial center
        const initialTransform = d3.zoomIdentity
            .translate(containerElement.clientWidth / 2, 100)
            .scale(0.8);
        
        d3.select(svgElement).call(zoom.transform, initialTransform);
    });

    export function reset() {
        const transform = d3.zoomIdentity.translate(containerElement.clientWidth / 2, 100).scale(0.8);
        d3.select(svgElement).transition().duration(750).call(zoom.transform, transform);
    }
</script>

<div bind:this={containerElement} class="pan-zoom-container" style="width: {width}; height: {height};">
    <svg bind:this={svgElement} class="pan-zoom-svg">
        <g bind:this={gElement}>
            {@render children()}
        </g>
    </svg>

    <div class="controls">
        <button onclick={reset} title="Reset View">Center</button>
    </div>
</div>

<style>
    .pan-zoom-container {
        position: relative;
        overflow: hidden;
        background: #f1f5f9;
        background-image: radial-gradient(#cbd5e1 1px, transparent 1px);
        background-size: 24px 24px;
    }

    .pan-zoom-svg {
        width: 100%;
        height: 100%;
        cursor: grab;
    }

    .pan-zoom-svg:active {
        cursor: grabbing;
    }

    .controls {
        position: absolute;
        bottom: var(--space-md);
        right: var(--space-md);
        display: flex;
        flex-direction: column;
        gap: var(--space-xs);
    }

    .controls button {
        padding: var(--space-xs) var(--space-sm);
        background: white;
        border: 1px solid var(--border-color);
        border-radius: var(--radius-sm);
        font-size: 0.75rem;
        font-weight: 600;
        cursor: pointer;
        box-shadow: var(--shadow-sm);
    }
</style>
