<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    Button,
    IconSpaceFixer,
  } from "@rilldata/web-common/components/button";
  import Forward from "@rilldata/web-common/components/icons/Forward.svelte";
  import Tooltip from "@rilldata/web-common/components/tooltip/Tooltip.svelte";
  import TooltipContent from "@rilldata/web-common/components/tooltip/TooltipContent.svelte";
  import { behaviourEvent } from "@rilldata/web-common/metrics/initMetrics";
  import { BehaviourEventMedium } from "@rilldata/web-common/metrics/service/BehaviourEventTypes";
  import {
    MetricsEventScreenName,
    MetricsEventSpace,
  } from "@rilldata/web-common/metrics/service/MetricsTypes";
  import { createRuntimeServiceGetFile } from "@rilldata/web-common/runtime-client";
  import { runtime } from "@rilldata/web-common/runtime-client/runtime-store";
  import { getFilePathFromNameAndType } from "../../entity-management/entity-mappers";
  import {
    fileArtifactsStore,
    getFileArtifactReconciliationErrors,
  } from "../../entity-management/file-artifacts-store";
  import { EntityType } from "../../entity-management/types";

  export let metricsDefName;

  $: fileQuery = createRuntimeServiceGetFile(
    $runtime.instanceId,
    getFilePathFromNameAndType(metricsDefName, EntityType.MetricsDefinition)
  );
  $: yaml = $fileQuery?.data?.blob;
  $: errors = getFileArtifactReconciliationErrors(
    $fileArtifactsStore,
    `${metricsDefName}.yaml`
  );

  let buttonDisabled = true;
  let buttonStatus;

  const viewDashboard = () => {
    goto(`/dashboard/${metricsDefName}`);

    behaviourEvent.fireNavigationEvent(
      metricsDefName,
      BehaviourEventMedium.Button,
      MetricsEventSpace.Workspace,
      MetricsEventScreenName.MetricsDefinition,
      MetricsEventScreenName.Dashboard
    );
  };

  const TOOLTIP_CTA = "Fix this error to enable your dashboard.";
  // no content
  $: if (!yaml?.length) {
    buttonDisabled = true;
    buttonStatus = [
      "Your metrics definition is empty. Get started by trying one of the options in the editor.",
    ];
  }
  // content & errors
  else if (errors?.length) {
    buttonDisabled = true;
    buttonStatus = [errors[0].message, TOOLTIP_CTA];
  }
  // preview is available
  else {
    buttonStatus = ["Explore your metrics dashboard"];
    buttonDisabled = false;
  }
</script>

<Tooltip alignment="middle" distance={5} location="right">
  <!-- TODO: we need to standardize these buttons. -->
  <Button
    disabled={buttonDisabled}
    label="Go to dashboard"
    on:click={() => viewDashboard()}
    type="primary"
  >
    <IconSpaceFixer pullLeft>
      <Forward /></IconSpaceFixer
    > Go to Dashboard
  </Button>
  <TooltipContent slot="tooltip-content">
    {#each buttonStatus as status}
      <div>{status}</div>
    {/each}
  </TooltipContent>
</Tooltip>
