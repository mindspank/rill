<script lang="ts">
  import { appQueryStatusStore } from "@rilldata/web-common/runtime-client/application-store";
  import { WorkspaceHeader } from "../../layout/workspace";

  export let fileName: string;
</script>

<div class="grid items-center" style:grid-template-columns="auto max-content">
  <WorkspaceHeader
    titleInput={fileName}
    onChangeCallback={undefined}
    appRunning={$appQueryStatusStore}
    editable={false}
    showInspectorToggle={false}
  />
</div>
