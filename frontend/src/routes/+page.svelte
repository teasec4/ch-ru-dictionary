<script lang="ts">
    let searchTerm = $state("");
    let responseData = $state<SearchResponse | null>(null);
    let loading = $state(false);
    let error = $state<string | null>(null);
    
    $effect(() => {
        console.log(responseData);
    })
    
    type Response = {
      chinese: string;
      pinyin: string;
      meanings: string[];
    }
    
    type SearchResponse = {
      data: Response[];
      count: number;
      message?: string;
    }

    async function handleClick() {
      if(searchTerm.trim()){
        loading = true;
        error = null;
        try {
          const response = await fetch("/api/search", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ word: searchTerm }),
          });
          if (response.ok) {
            const data = await response.json();
            searchTerm = "";
            responseData = data;
          } else {
            responseData = null;
          }
        } catch (e) {
          error = e instanceof Error ? e.message : String(e);
        } finally {
          loading = false;
        }
      }
    }
</script>

<div class="w-screen h-screen flex flex-col">
    <div class="flex gap-3 py-8 flex-col justify-center items-center">
        <input 
            type="text" 
            placeholder="" 
            bind:value={searchTerm} 
            class="px-2 py-2 text-2xl max-w-sm w-90 rounded-md border border-text-primary focus:outline-1 focus:ring-1 focus:ring-text-primary"
        />
        <button 
            onclick={handleClick} 
            disabled={loading && !error && !searchTerm.trim()}
            class="px-6 py-2 rounded-md bg-text-primary text-text-tertiary font-semibold disabled:bg-gray-400"
        >
            {#if !loading}Поиск{/if}
        </button>
    </div>
   
    <div class="flex-1 flex flex-col justify-start items-center p-4">
        {#if responseData}
            {#each responseData.data as word}
                <div class="w-full max-w-md flex flex-col gap-2">
                    <div class="text-2xl font-bold">{word.chinese}</div>
                    <div class="text-gray-700 italic">{word.pinyin}</div>
                    <div class="text-gray-800">Значения: {word.meanings}</div>
                </div>
            {/each}
            {#if responseData.message === "no results"}
                <p class="text-2xl text-text-primary mt-4">Нет результатов</p>
            {/if}
        {/if}
       
        {#if !responseData && !loading && !error}
            <p class="text-text-primary text-2xl mt-4">Введите слово для поиска</p>
        {/if}
        
        {#if error}
        <p class="text-red-600 mb-4">Ошибка: {error}</p>
        {/if}
    </div>
    
</div>

<style lang="postcss">
    @reference "tailwindcss";

</style>