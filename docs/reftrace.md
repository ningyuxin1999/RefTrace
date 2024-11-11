<div class="pdoc" role="main">

<div class="section module-info">

# reftrace

</div>

<div id="Module" class="section">

<div class="attr class">

<span class="def">class</span> <span class="name">Module</span>:

</div>

<a href="#Module" class="headerlink"></a>

<div id="Module.__init__" class="classattr">

<div class="attr function">

<span class="name">Module</span><span class="signature pdoc-code condensed">(<span class="param"><span class="n">filepath</span><span class="p">:</span>
<span class="nb">str</span></span>)</span>

</div>

<a href="#Module.__init__" class="headerlink"></a>

</div>

<div id="Module.path" class="classattr">

<div class="attr variable">

<span class="name">path</span><span class="annotation">: str</span>

</div>

<a href="#Module.path" class="headerlink"></a>

</div>

<div id="Module.dsl_version" class="classattr">

<div class="attr variable">

<span class="name">dsl_version</span><span class="annotation">:
int</span>

</div>

<a href="#Module.dsl_version" class="headerlink"></a>

</div>

<div id="Module.processes" class="classattr">

<div class="attr variable">

<span class="name">processes</span><span class="annotation">:
list\[[Process](#Process)\]</span>

</div>

<a href="#Module.processes" class="headerlink"></a>

</div>

</div>

<div id="Process" class="section">

<div class="attr class">

<span class="def">class</span> <span class="name">Process</span>:

</div>

<a href="#Process" class="headerlink"></a>

<div id="Process.__init__" class="classattr">

<div class="attr function">

<span class="name">Process</span><span class="signature pdoc-code condensed">(<span class="param"><span class="n">handle</span><span class="p">:</span>
<span class="nb">int</span></span>)</span>

</div>

<a href="#Process.__init__" class="headerlink"></a>

</div>

<div id="Process.name" class="classattr">

<div class="attr variable">

<span class="name">name</span><span class="annotation">: str</span>

</div>

<a href="#Process.name" class="headerlink"></a>

</div>

<div id="Process.containers" class="classattr">

<div class="attr variable">

<span class="name">containers</span><span class="annotation">:
list\[[Container](#Container)\]</span>

</div>

<a href="#Process.containers" class="headerlink"></a>

</div>

<div id="Process.labels" class="classattr">

<div class="attr variable">

<span class="name">labels</span>

</div>

<a href="#Process.labels" class="headerlink"></a>

</div>

</div>

<div id="Container" class="section">

<div class="attr class">

<span class="def">class</span> <span class="name">Container</span>:

</div>

<a href="#Container" class="headerlink"></a>

<div id="Container.__init__" class="classattr">

<div class="attr function">

<span class="name">Container</span><span class="signature pdoc-code condensed">(<span class="param"><span class="n">handle</span><span class="p">:</span>
<span class="nb">int</span></span>)</span>

</div>

<a href="#Container.__init__" class="headerlink"></a>

</div>

<div id="Container.format" class="classattr">

<div class="attr variable">

<span class="name">format</span><span class="annotation">: str</span>

</div>

<a href="#Container.format" class="headerlink"></a>

</div>

<div id="Container.simple_name" class="classattr">

<div class="attr variable">

<span class="name">simple_name</span><span class="annotation">:
Optional\[str\]</span>

</div>

<a href="#Container.simple_name" class="headerlink"></a>

</div>

<div id="Container.condition" class="classattr">

<div class="attr variable">

<span class="name">condition</span><span class="annotation">:
Optional\[str\]</span>

</div>

<a href="#Container.condition" class="headerlink"></a>

</div>

<div id="Container.true_name" class="classattr">

<div class="attr variable">

<span class="name">true_name</span><span class="annotation">:
Optional\[str\]</span>

</div>

<a href="#Container.true_name" class="headerlink"></a>

</div>

<div id="Container.false_name" class="classattr">

<div class="attr variable">

<span class="name">false_name</span><span class="annotation">:
Optional\[str\]</span>

</div>

<a href="#Container.false_name" class="headerlink"></a>

</div>

</div>

<div id="Label" class="section">

<div class="attr class">

<span class="def">class</span> <span class="name">Label</span>:

</div>

<a href="#Label" class="headerlink"></a>

<div id="Label.__init__" class="classattr">

<div class="attr function">

<span class="name">Label</span><span class="signature pdoc-code condensed">(<span class="param"><span class="n">handle</span><span class="p">:</span>
<span class="nb">int</span></span>)</span>

</div>

<a href="#Label.__init__" class="headerlink"></a>

</div>

<div id="Label.value" class="classattr">

<div class="attr variable">

<span class="name">value</span><span class="annotation">: str</span>

</div>

<a href="#Label.value" class="headerlink"></a>

</div>

</div>

</div>
