<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<h1>House</h1>

<div>
  <p>
    Location : <strong>{data.house.location}</strong>, Block {data.house.block},
    Partition {data.house.partition}, Occuppied {data.house.occupied}
  </p>
</div>
