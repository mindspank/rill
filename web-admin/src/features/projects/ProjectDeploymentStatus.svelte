<script lang="ts">
  import { goto } from "$app/navigation";
  import Button from "@rilldata/web-common/components/button/Button.svelte";
  import { createAdminServiceGetProject } from "../../client";
  import ProjectDeploymentStatusChip from "./ProjectDeploymentStatusChip.svelte";

  export let organization: string;
  export let project: string;

  $: proj = createAdminServiceGetProject(organization, project);
  $: isProjectDeployed = $proj?.data && $proj.data.prodDeployment;

  function handleViewLogs() {
    goto(`/${organization}/${project}/-/logs`);
  }
</script>

<div class="flex flex-col gap-y-2">
  <div>
    <span class="uppercase text-gray-500 font-semibold text-[10px] leading-4"
      >Project status</span
    >
    <div>
      <ProjectDeploymentStatusChip {organization} {project} />
    </div>
    {#if $proj && $proj.data && $proj.data.prodDeployment}
      <span class="text-gray-500 text-[11px] leading-4">
        Synced {new Date($proj.data.prodDeployment.updatedOn).toLocaleString(
          undefined,
          {
            month: "short",
            day: "numeric",
            hour: "numeric",
            minute: "numeric",
          }
        )}
      </span>
    {/if}
  </div>
  {#if isProjectDeployed}
    <div>
      <Button type="secondary" on:click={handleViewLogs}>View logs</Button>
    </div>
  {/if}
</div>
