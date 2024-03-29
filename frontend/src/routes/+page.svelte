<script lang="ts">
  import type { ActionData } from "./$types";
  export let form: ActionData;
</script>

<h3>Login</h3>
<form method="POST">
  <label>
    Phone
    <input name="phone" type="phone" />
  </label>
  <label>
    Password
    <input name="password" type="password" />
  </label>
  <button>Log in</button>
</form>

{#if form?.success === false}
  <div class="alert alert-error">
    <div>
      {form.error}
    </div>
  </div>
{/if}
