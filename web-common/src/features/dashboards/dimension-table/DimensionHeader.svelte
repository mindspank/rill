<script lang="ts">
  import { Switch } from "@rilldata/web-common/components/button";
  import Back from "@rilldata/web-common/components/icons/Back.svelte";
  import Close from "@rilldata/web-common/components/icons/Close.svelte";
  import SearchIcon from "@rilldata/web-common/components/icons/Search.svelte";
  import { Search } from "@rilldata/web-common/components/search";
  import Shortcut from "@rilldata/web-common/components/tooltip/Shortcut.svelte";
  import Tooltip from "@rilldata/web-common/components/tooltip/Tooltip.svelte";
  import TooltipContent from "@rilldata/web-common/components/tooltip/TooltipContent.svelte";
  import TooltipShortcutContainer from "@rilldata/web-common/components/tooltip/TooltipShortcutContainer.svelte";
  import TooltipTitle from "@rilldata/web-common/components/tooltip/TooltipTitle.svelte";
  import { cancelDashboardQueries } from "@rilldata/web-common/features/dashboards/dashboard-queries";
  import { EntityStatus } from "@rilldata/web-common/features/entity-management/types";
  import { slideRight } from "@rilldata/web-common/lib/transitions";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { createEventDispatcher } from "svelte";
  import { fly } from "svelte/transition";
  import Spinner from "../../entity-management/Spinner.svelte";
  import { metricsExplorerStore } from "../dashboard-stores";
  import ExportDimensionTableDataButton from "./ExportDimensionTableDataButton.svelte";

  export let metricViewName: string;
  export let dimensionName: string;
  export let isFetching: boolean;
  export let excludeMode = false;

  const queryClient = useQueryClient();

  $: filterKey = excludeMode ? "exclude" : "include";
  $: otherFilterKey = excludeMode ? "include" : "exclude";

  let searchToggle = false;

  const dispatch = createEventDispatcher();

  let searchText = "";
  function onSearch() {
    dispatch("search", searchText);
  }

  function closeSearchBar() {
    searchText = "";
    searchToggle = !searchToggle;
    onSearch();
  }

  const goBackToLeaderboard = () => {
    metricsExplorerStore.setMetricDimensionName(metricViewName, null);
  };
  function toggleFilterMode() {
    cancelDashboardQueries(queryClient, metricViewName);
    metricsExplorerStore.toggleFilterMode(metricViewName, dimensionName);
  }
</script>

<div class="flex justify-between items-center p-1">
  <button class="flex items-center" on:click={() => goBackToLeaderboard()}>
    {#if isFetching}
      <div>
        <Spinner size="16px" status={EntityStatus.Running} />
      </div>
    {:else}
      <span class="ui-copy-icon">
        <Back size="16px" />
      </span>
      <span> All Dimensions </span>
    {/if}
  </button>

  <!-- We fix the height to avoid a layout shift when the Search component is expanded. -->
  <div class="flex items-center gap-x-5 cursor-pointer h-9">
    <Tooltip distance={16} location="left">
      <div class="flex items-center gap-x-1 ui-copy-icon">
        <Switch checked={excludeMode} on:click={() => toggleFilterMode()}>
          Exclude
        </Switch>
      </div>
      <TooltipContent slot="tooltip-content">
        <TooltipTitle>
          <svelte:fragment slot="name">
            Output {filterKey}s selected values
          </svelte:fragment>
        </TooltipTitle>
        <TooltipShortcutContainer>
          <div>Toggle to {otherFilterKey} values</div>
          <Shortcut>Click</Shortcut>
        </TooltipShortcutContainer>
      </TooltipContent>
    </Tooltip>

    {#if !searchToggle}
      <button
        class="flex items-center gap-x-1 text-gray-700"
        in:fly={{ x: 10, duration: 300 }}
        on:click={() => (searchToggle = !searchToggle)}
      >
        <SearchIcon size="16px" />
        <span>Search</span>
      </button>
    {:else}
      <div
        transition:slideRight|local={{ leftOffset: 8 }}
        class="flex items-center gap-x-1"
      >
        <Search bind:value={searchText} on:input={onSearch} />
        <button class="ui-copy-icon" on:click={() => closeSearchBar()}>
          <Close />
        </button>
      </div>
    {/if}

    <ExportDimensionTableDataButton {metricViewName} />
  </div>
</div>
