<script lang="ts">
  import * as d3 from "d3";
  import PersonCard from "./PersonCard.svelte";

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
    accessType = "view",
    onEdit,
    onDelete,
    onAddFamily,
  } = $props();

  // Layout constants
  const nodeWidth = 180;
  const nodeHeight = 140;
  const hGap = 60; // horizontal gap between nodes
  const vGap = 120; // vertical gap between levels
  const spouseGap = 40; // gap between spouses

  // Computed
  let positions = $state<Map<string, { x: number; y: number }>>(new Map());
  let parentLinks = $state<any[]>([]);
  let spouseLines = $state<any[]>([]);

  function buildLayout() {
    positions = new Map();
    parentLinks = [];
    spouseLines = [];

    if (people.length === 0) return;

    const parentChildRels = relationships.filter(
      (r) => r.type === "parent_child",
    );
    const spouseRels = relationships.filter((r) => r.type === "spouse");

    // Children set (person_b_id in parent_child)
    const childrenSet = new Set(parentChildRels.map((r) => r.person_b_id));

    // Roots: people not a child in any parent_child rel
    const roots = people.filter((p) => !childrenSet.has(p.id));

    // Build a spouse-pair map  id -> partnerId
    const spouseMap = new Map<string, string>();
    for (const r of spouseRels) {
      spouseMap.set(r.person_a_id, r.person_b_id);
      spouseMap.set(r.person_b_id, r.person_a_id);
    }

    // Choose canonical representative for each couple
    // (person_a_id is canonical to avoid duplicating layout)
    const handledRoots = new Set<string>();
    const rootGroups: string[][] = []; // each group is [personId] or [personA, personB] for couples

    for (const root of roots) {
      if (handledRoots.has(root.id)) continue;
      const partnerId = spouseMap.get(root.id);
      // only group if partner is also a root (not a child of someone)
      if (
        partnerId &&
        roots.some((r) => r.id === partnerId) &&
        !handledRoots.has(partnerId)
      ) {
        // canonical: use person_a_id from spouseRels
        const rel = spouseRels.find(
          (r) =>
            (r.person_a_id === root.id && r.person_b_id === partnerId) ||
            (r.person_a_id === partnerId && r.person_b_id === root.id),
        );
        const a = rel ? rel.person_a_id : root.id;
        const b = rel ? rel.person_b_id : partnerId;
        rootGroups.push([a, b]);
        handledRoots.add(a);
        handledRoots.add(b);
      } else if (partnerId && childrenSet.has(partnerId)) {
        // This root's spouse is a child — skip now, layoutNode will place them together
        // Don't add to rootGroups; don't mark as handled so it's picked up during child layout
      } else {
        // Lone root with no spouse or spouse is a separate root not yet seen
        rootGroups.push([root.id]);
        handledRoots.add(root.id);
      }
    }

    // Collect all unique children recursively for each root group
    const childrenOf = (
      personId: string,
      visited = new Set<string>(),
    ): string[] => {
      const direct = parentChildRels
        .filter((r) => r.person_a_id === personId)
        .map((r) => r.person_b_id);

      // Include partner's children too so the group subtree width is computed
      const partnerId = spouseMap.get(personId);
      const partnerDirect = partnerId
        ? parentChildRels
            .filter((r) => r.person_a_id === partnerId)
            .map((r) => r.person_b_id)
        : [];

      const all = [...new Set([...direct, ...partnerDirect])].filter(
        (id) => !visited.has(id),
      );
      return all;
    };

    // Layout algorithm: DFS placing nodes left-to-right
    let cursor = 0;

    function layoutGroup(group: string[], y: number): number {
      // group is [personId] or [personA, personB]
      // Returns the center x of this group

      // Collect unique children for the whole group
      const childIds = new Set<string>();
      for (const id of group) {
        for (const cid of childrenOf(id)) {
          if (!positions.has(cid)) childIds.add(cid);
        }
      }

      // Recursively lay out children first to know their span
      const childCenters: number[] = [];
      for (const cid of childIds) {
        const childPerson = people.find((p) => p.id === cid);
        if (!childPerson) continue;
        const cc = layoutNode(cid, y + nodeHeight + vGap);
        childCenters.push(cc);
      }

      // Center of parent group = center of children range (or cursor if leaf)
      let groupCenterX: number;
      if (childCenters.length > 0) {
        groupCenterX =
          (Math.min(...childCenters) + Math.max(...childCenters)) / 2;
      } else {
        // leaf: place at cursor
        const halfWidth =
          group.length === 2 ? nodeWidth + spouseGap / 2 : nodeWidth / 2;
        groupCenterX = cursor + halfWidth;
        cursor +=
          group.length === 2
            ? nodeWidth * 2 + spouseGap + hGap
            : nodeWidth + hGap;
      }

      // Place group members
      if (group.length === 2) {
        const [a, b] = group;
        positions.set(a, {
          x: groupCenterX - nodeWidth / 2 - spouseGap / 2,
          y,
        });
        positions.set(b, {
          x: groupCenterX + nodeWidth / 2 + spouseGap / 2,
          y,
        });
      } else {
        positions.set(group[0], { x: groupCenterX, y });
      }

      // Guard cursor to not go backwards
      const rightEdge =
        group.length === 2
          ? groupCenterX + nodeWidth + spouseGap / 2 + hGap
          : groupCenterX + nodeWidth / 2 + hGap;
      if (rightEdge > cursor) cursor = rightEdge;

      return groupCenterX;
    }

    function layoutNode(personId: string, y: number): number {
      if (positions.has(personId)) {
        return positions.get(personId)!.x;
      }

      const partnerId = spouseMap.get(personId);
      // If this child has an unplaced spouse, place them side by side as a couple
      if (partnerId && !positions.has(partnerId)) {
        // Place as a couple: child on left, spouse on right
        const coupleCenter = cursor + nodeWidth + spouseGap / 2;
        positions.set(personId, { x: coupleCenter - nodeWidth / 2 - spouseGap / 2, y });
        positions.set(partnerId, { x: coupleCenter + nodeWidth / 2 + spouseGap / 2, y });
        cursor += nodeWidth * 2 + spouseGap + hGap;

        // Recursively lay out their children at the next level
        const childIds = new Set<string>();
        for (const id of [personId, partnerId]) {
          for (const cid of childrenOf(id)) {
            if (!positions.has(cid)) childIds.add(cid);
          }
        }
        const childCenters: number[] = [];
        for (const cid of childIds) {
          const cc = layoutNode(cid, y + nodeHeight + vGap);
          childCenters.push(cc);
        }
        // Re-center the couple over their children if they have any
        if (childCenters.length > 0) {
          const newCenter = (Math.min(...childCenters) + Math.max(...childCenters)) / 2;
          positions.set(personId, { x: newCenter - nodeWidth / 2 - spouseGap / 2, y });
          positions.set(partnerId, { x: newCenter + nodeWidth / 2 + spouseGap / 2, y });
        }

        return positions.get(personId)!.x;
      }

      // Single node placement
      positions.set(personId, { x: cursor + nodeWidth / 2, y });
      cursor += nodeWidth + hGap;

      return positions.get(personId)!.x;
    }

    // Layout root groups first
    for (const group of rootGroups) {
      layoutGroup(group, 0);
    }

    // Place any remaining people who have no position yet
    // If they are a spouse of an already-placed person, match that person's Y level
    for (const p of people) {
      if (!positions.has(p.id)) {
        const partnerId = spouseMap.get(p.id);
        const partnerPos = partnerId ? positions.get(partnerId) : undefined;
        const y = partnerPos ? partnerPos.y : 0;
        positions.set(p.id, { x: cursor + nodeWidth / 2, y });
        cursor += nodeWidth + hGap;
      }
    }

    // Build parent-child link lines (grouped by parent-pair/single-parent)
    const parentUnits = new Map<string, string[]>(); // childId -> [parentIds]
    for (const rel of parentChildRels) {
      if (!parentUnits.has(rel.person_b_id))
        parentUnits.set(rel.person_b_id, []);
      parentUnits.get(rel.person_b_id)!.push(rel.person_a_id);
    }

    for (const [childId, parentIds] of parentUnits.entries()) {
      const tgtPos = positions.get(childId);
      if (!tgtPos) continue;

      // Check if any parent has a spouse (even if only one parent_child record exists)
      let srcX: number | null = null;
      let srcY: number | null = null;

      // Collect all relevant parents: those explicitly listed + their spouses
      const allParentIds = new Set<string>(parentIds);
      for (const pid of parentIds) {
        const spouseId = spouseMap.get(pid);
        if (spouseId) allParentIds.add(spouseId);
      }

      // Find a spouse pair among all parents
      let foundPair: [string, string] | null = null;
      for (const pid of parentIds) {
        const spouseId = spouseMap.get(pid);
        if (spouseId && allParentIds.has(spouseId)) {
          foundPair = [pid, spouseId];
          break;
        }
        // Also check if the parent has a spouse even if that spouse isn't in parentIds
        if (spouseId) {
          const spousePos = positions.get(spouseId);
          if (spousePos) {
            foundPair = [pid, spouseId];
            break;
          }
        }
      }

      if (foundPair) {
        const [p1, p2] = foundPair;
        const pos1 = positions.get(p1);
        const pos2 = positions.get(p2);
        if (pos1 && pos2) {
          // Start from center of the spouse line
          srcX = (pos1.x + pos2.x) / 2;
          srcY = (pos1.y + pos2.y) / 2;
        }
      }

      if (srcX === null) {
        // Single parent with no spouse: draw from that parent's center
        for (const pid of parentIds) {
          const pos = positions.get(pid);
          if (pos) {
            parentLinks.push({
              src: { x: pos.x, y: pos.y },
              tgt: { x: tgtPos.x, y: tgtPos.y - nodeHeight / 2 },
            });
          }
        }
      } else {
        // Draw a single symmetric line from the couple's midpoint
        parentLinks.push({
          src: { x: srcX, y: srcY },
          tgt: { x: tgtPos.x, y: tgtPos.y - nodeHeight / 2 },
        });
      }
    }

    // Build spouse lines
    for (const rel of spouseRels) {
      const posA = positions.get(rel.person_a_id);
      const posB = positions.get(rel.person_b_id);
      if (posA && posB) {
        spouseLines.push({ a: posA, b: posB });
      }
    }
  }

  $effect(() => {
    buildLayout();
  });

  // Curved path with a vertical "stem" that exits the parent card area
  function parentLinkPath(
    src: { x: number; y: number },
    tgt: { x: number; y: number },
  ): string {
    const sx = src.x;
    const sy = src.y;
    const tx = tgt.x;
    const ty = tgt.y;

    // Point where the line should exit the parent card area (below bottom edge)
    const stemY = sy + nodeHeight / 2 + 30;

    const path = d3.path();
    path.moveTo(sx, sy);
    path.lineTo(sx, stemY); // Vertical stem down through the gap

    const my = (stemY + ty) / 2;
    path.bezierCurveTo(sx, my, tx, my, tx, ty); // Curve from stem tip to child top
    return path.toString();
  }

  // Horizontal spouse connector between right-edge of A and left-edge of B
  function spouseLinkPath(
    a: { x: number; y: number },
    b: { x: number; y: number },
  ): string {
    // Make sure left is always 'a'
    const left = a.x < b.x ? a : b;
    const right = a.x < b.x ? b : a;
    const leftX = left.x + nodeWidth / 2;
    const rightX = right.x - nodeWidth / 2;
    const midY = (left.y + right.y) / 2;

    if (Math.abs(left.y - right.y) < 5) {
      // Same-level: simple horizontal line at card center
      const path = d3.path();
      path.moveTo(leftX, midY);
      path.lineTo(rightX, midY);
      return path.toString();
    } else {
      // Different levels: S-curve
      const midX = (leftX + rightX) / 2;
      const path = d3.path();
      path.moveTo(leftX, left.y);
      path.bezierCurveTo(midX, left.y, midX, right.y, rightX, right.y);
      return path.toString();
    }
  }

  // Compute SVG viewBox from positions
  let svgPadding = 80;
  let svgWidth = $derived.by(() => {
    if (positions.size === 0) return 600;
    const xs = [...positions.values()].map((p) => p.x);
    return Math.max(...xs) + nodeWidth / 2 + svgPadding * 2;
  });
  let svgHeight = $derived.by(() => {
    if (positions.size === 0) return 400;
    const ys = [...positions.values()].map((p) => p.y);
    return Math.max(...ys) + nodeHeight + 60 + svgPadding * 2;
  });
</script>

<g class="tree-graph">
  <!-- Spouse connector lines (drawn FIRST so they appear behind nodes) -->
  <g class="spouse-links">
    {#each spouseLines as link}
      <path
        d={spouseLinkPath(link.a, link.b)}
        fill="none"
        stroke="#f43f5e"
        stroke-width="3"
        stroke-dasharray="8,5"
        stroke-linecap="round"
        class="spouse-connector"
      />
    {/each}
  </g>

  <!-- Parent-child connector lines -->
  <g class="parent-links">
    {#each parentLinks as link}
      <path
        d={parentLinkPath(link.src, link.tgt)}
        fill="none"
        stroke="#cbd5e1"
        stroke-width="2.5"
        stroke-linecap="round"
        class="parent-connector"
      />
    {/each}
  </g>

  <!-- Person card nodes -->
  <g class="nodes">
    {#each [...positions.entries()] as [personId, pos]}
      {@const person = people.find((p) => p.id === personId)}
      {#if person}
        <foreignObject
          x={pos.x - nodeWidth / 2}
          y={pos.y - nodeHeight / 2}
          width={nodeWidth}
          height={nodeHeight + 50}
          class="node-fo"
        >
          <PersonCard {person} {accessType} {onEdit} {onDelete} {onAddFamily} />
        </foreignObject>
      {/if}
    {/each}
  </g>
</g>

<style>
  .spouse-connector {
    opacity: 0.85;
  }

  .parent-connector {
    transition: d 0.3s ease;
  }

  .node-fo {
    overflow: visible;
  }
</style>
