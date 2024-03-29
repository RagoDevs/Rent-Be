<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<h1>Welcome to Houses</h1>

{#each data.houses as house}
  <div>
    <a href={`/auth/houses/${house.id}`}
      ><p>
        Location : <strong>{house.location}</strong>, Block {house.block},
        Partition {house.partition}, Occuppied {house.occupied}
      </p></a
    >
  </div>
{/each}
