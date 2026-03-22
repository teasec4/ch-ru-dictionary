<script lang="ts">
    let searchTerm = $state("");
    let responseData = $state<Response | null>(null);
    
    $effect(() => {
        console.log(responseData);
    })
    
    type Response = {
      chinese: string;
      pinyin: string;
      pinyin_normalized: string;
      meanings: string[];
    }

    async function handleClick() {
      if(searchTerm.trim()){
        let response = await fetch(`/api/search/${searchTerm}`); 
        if (response.ok) {
          let data = await response.json();
          responseData = data;
        } else {
          responseData = null;
        }
      }
    }
</script>

<h1>Словарь</h1>
<input type="text" placeholder="Search..." bind:value={searchTerm} />
<button onclick={handleClick}>поиск</button>

{#if responseData}
  <div>
    <span>{responseData.chinese}</span>
    <span>{responseData.pinyin}</span>
    <span>{responseData.meanings.join(", ")}</span>
  </div>
{/if}

{#if responseData === null}
  <p>Пусто</p>
{/if}

<style>
    div {
        display: flex;
        flex-direction: column;
        align-items: start;
    }
</style>