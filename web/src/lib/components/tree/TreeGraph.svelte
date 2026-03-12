<script lang="ts">
    import { onMount } from 'svelte';
    import * as d3 from 'd3';
    import PersonCard from './PersonCard.svelte';

    /**
     * @typedef {Object} Props
     * @property {any[]} people
     * @property {any[]} relationships
     * @property {string} accessType
     * @property {Function} onEdit
     * @property {Function} onDelete
     * @property {Function} onAddFamily
     */
    let { 
        people = [], 
        relationships = [], 
        accessType = 'view',
        onEdit,
        onDelete,
        onAddFamily 
    } = $props();

    // D3 Layout variables
    const nodeWidth = 200;
    const nodeHeight = 120;
    const horizontalGap = 40;
    const verticalGap = 100;

    // Computed tree data
    let nodes = $state<any[]>([]);
    let links = $state<any[]>([]);

    function buildLayout() {
        if (people.length === 0) {
            nodes = [];
            links = [];
            return;
        }

        // 1. Identify parents and children from relationships
        const childrenIds = new Set(
            relationships
                .filter(rel => rel.Type === 'parent_child')
                .map(rel => rel.PersonBID)
        );

        // 2. Find root nodes (people who are not children in any parent_child relationship)
        const roots = people.filter(p => !childrenIds.has(p.ID));

        // 3. Create hierarchy
        // Note: Silsilah can be complex (multiple roots, multiple parents), 
        // but d3-hierarchy expects a single root. For now, we'll create a virtual root.
        const data = {
            id: 'virtual-root',
            name: 'Virtual Root',
            children: roots.map(root => {
                return buildSubtree(root, people, relationships);
            })
        };

        const hierarchy = d3.hierarchy(data);
        const treeLayout = d3.tree<any>().nodeSize([nodeWidth + horizontalGap, nodeHeight + verticalGap]);
        const treeData = treeLayout(hierarchy);

        // Filter out virtual root
        nodes = treeData.descendants().filter((d: any) => d.data.id !== 'virtual-root');
        links = treeData.links().filter((l: any) => l.source.data.id !== 'virtual-root');
    }

    function buildSubtree(person: any, allPeople: any[], allRels: any[]) {
        const item: any = { ...person };
        
        // Find children
        const childRels = allRels.filter(rel => rel.Type === 'parent_child' && rel.PersonAID === person.ID);
        if (childRels.length > 0) {
            item.children = childRels.map(rel => {
                const child = allPeople.find(p => p.ID === rel.PersonBID);
                return child ? buildSubtree(child, allPeople, allRels) : null;
            }).filter(Boolean);
        }
        
        return item;
    }

    $effect(() => {
        buildLayout();
    });

    // Helper for curved path
    function generateLinkPath(link: any) {
        const sourceX = link.source.x;
        const sourceY = link.source.y + nodeHeight / 2;
        const targetX = link.target.x;
        const targetY = link.target.y - nodeHeight / 2;

        const deltaY = targetY - sourceY;
        
        const path = d3.path();
        path.moveTo(sourceX, sourceY);
        path.bezierCurveTo(
            sourceX, sourceY + deltaY / 2,
            targetX, targetY - deltaY / 2,
            targetX, targetY
        );
        return path.toString();
    }
</script>

<g class="tree-graph">
    <!-- Links Layer -->
    <g class="links">
        {#each links as link}
            <path
                d={generateLinkPath(link)}
                class="connector"
                fill="none"
                stroke="#cbd5e1"
                stroke-width="2"
            />
        {/each}
    </g>

    <!-- Nodes Layer -->
    <g class="nodes">
        {#each nodes as node}
            <foreignObject
                x={node.x - nodeWidth / 2}
                y={node.y - nodeHeight / 2}
                width={nodeWidth}
                height={nodeHeight + 40}
                class="node-container"
            >
                <PersonCard
                    person={node.data}
                    {accessType}
                    {onEdit}
                    {onDelete}
                    {onAddFamily}
                />
            </foreignObject>
        {/each}
    </g>
</g>

<style>
    .connector {
        transition: d 0.3s ease;
    }

    .node-container {
        overflow: visible;
    }
</style>
