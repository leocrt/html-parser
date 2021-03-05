package main

import (
	"strings"

	parser "html-parse/parser"
)

const CHAPTER_FONT = "{font-size:18px;font-family:Times New Roman,Bold;color:#000000;}"
const ARTICLE_FONT = "{font-size:18px;font-family:Times New Roman;color:#000000;}"
const LINE_BREAK_FONT = "{font-size:18px;line-height:20px;font-family:Times New Roman;color:#000000;}"

var exampleHtml = `
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" lang="" xml:lang="">
<head>
<title>Page 1</title>

<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
<style type="text/css">
<!--
	p {margin: 0; padding: 0;}	
	.ft00{font-size:16px;font-family:ABCDEE+Calibri;color:#000000;}
	.ft01{font-size:18px;font-family:ABCDEE+Calibri,Bold;color:#000000;}
	.ft02{font-size:18px;font-family:ABCDEE+Calibri;color:#933634;}
	.ft03{font-size:18px;font-family:ABCDEE+Calibri;color:#000000;}
	.ft04{font-size:17px;font-family:Times New Roman,Bold;color:#000000;}
	.ft05{font-size:18px;font-family:Times New Roman,Bold;color:#0000ff;}
	.ft06{font-size:18px;font-family:Times New Roman;color:#000000;}
	.ft07{font-size:18px;font-family:Times New Roman,Bold;color:#000000;}
	.ft08{font-size:18px;line-height:20px;font-family:Times New Roman;color:#000000;}
-->
</style>
</head>
<body bgcolor="#A0A0A0" vlink="blue" link="blue">
<div id="page1-div" style="position:relative;width:892px;height:1262px;">
<img width="892" height="1262" src="rdc301001.png" alt="background image"/>
<p style="position:absolute;top:117px;left:485px;white-space:nowrap" class="ft00">;</p>
<p style="position:absolute;top:133px;left:353px;white-space:nowrap" class="ft01"><b>Ministério;da;Saúde;-;MS;</b></p>
<p style="position:absolute;top:155px;left:260px;white-space:nowrap" class="ft01"><b>Agência;Nacional;de;Vigilância;Sanitária;–;ANVISA;</b></p>
<p style="position:absolute;top:177px;left:128px;white-space:nowrap" class="ft00">;</p>
<p style="position:absolute;top:1192px;left:198px;white-space:nowrap" class="ft02">Este;texto;não;substitui;o(s);publicado(s);em;Diário;Oficial;da;União.</p>
<p style="position:absolute;top:1192px;left:695px;white-space:nowrap" class="ft03">;</p>
<p style="position:absolute;top:198px;left:92px;white-space:nowrap" class="ft04"><b>RESOLUÇÃO;DA;DIRETORIA;COLEGIADA;-;RDC;Nº;301, DE;21;DE AGOSTO DE 2019;(*);</b></p>
<p style="position:absolute;top:233px;left:248px;white-space:nowrap" class="ft05"><b>(Publicada no DOU nº 162, de;22;de;agosto de 2019);</b></p>
<p style="position:absolute;top:269px;left:242px;white-space:nowrap" class="ft05"><b>(Republicada no DOU;nº 49, de;12 de;março de;2020);</b></p>
<p style="position:absolute;top:304px;left:248px;white-space:nowrap" class="ft05"><b>(Republicada no DOU;nº;78, de;24 de;abril;de;2020);</b></p>
<p style="position:absolute;top:340px;left:510px;white-space:nowrap" class="ft08">Dispõe sobre as Diretrizes Gerais;<br/>de Boas Práticas de Fabricação de;<br/>Medicamentos.;</p>
<p style="position:absolute;top:417px;left:170px;white-space:nowrap" class="ft06">A;</p>
<p style="position:absolute;top:417px;left:188px;white-space:nowrap" class="ft07"><b>Diretoria;Colegiada;da;Agência;Nacional;de;Vigilância;Sanitária</b></p>
<p style="position:absolute;top:417px;left:686px;white-space:nowrap" class="ft06">,;no;uso;da;</p>
<p style="position:absolute;top:437px;left:128px;white-space:nowrap" class="ft08">atribuição;que;lhe;confere;o;art.;15,;III;e;IV,;aliado;ao;art.;7º,;III;e;IV;da;Lei;nº;9.782,;<br/>de;26;de;janeiro;de;1999,;e;ao;art.;53,;V,;§§;1º;e;3º;do;Regimento;Interno;aprovado;pela;<br/>Resolução;da;Diretoria;Colegiada;–;RDC;n°;255,;de;10;de;dezembro;de;2018,;resolve;<br/>adotar a seguinte Resolução da Diretoria Colegiada, conforme deliberado em reunião;<br/>realizada;em;20;de;agosto;de;2019,;e;eu,;Diretor-Presidente;Substituto,;determino;a;sua;<br/>publicação.;</p>
<p style="position:absolute;top:577px;left:393px;white-space:nowrap" class="ft07"><b>CAPÍTULO;I;</b></p>
<p style="position:absolute;top:613px;left:325px;white-space:nowrap" class="ft07"><b>DAS;DISPOSIÇÕES;INICIAIS;</b></p>
<p style="position:absolute;top:649px;left:419px;white-space:nowrap" class="ft07"><b>Seção I;</b></p>
<p style="position:absolute;top:684px;left:402px;white-space:nowrap" class="ft07"><b>Do objetivo;</b></p>
<p style="position:absolute;top:719px;left:170px;white-space:nowrap" class="ft06">Art. 1º Esta Resolução possui o objetivo de adotar as diretrizes gerais de Boas;</p>
<p style="position:absolute;top:740px;left:128px;white-space:nowrap" class="ft08">Práticas de Fabricação de Medicamentos do Esquema de Cooperação em Inspeção;<br/>Farmacêutica, PIC/S, como requisitos mínimos a serem seguidos na fabricação de;<br/>medicamentos.;</p>
<p style="position:absolute;top:818px;left:415px;white-space:nowrap" class="ft07"><b>Seção II;</b></p>
<p style="position:absolute;top:853px;left:386px;white-space:nowrap" class="ft07"><b>Da abrangência;</b></p>
<p style="position:absolute;top:889px;left:170px;white-space:nowrap" class="ft06">Art. 2º Esta Resolução se aplica às empresas que realizam as operações;</p>
<p style="position:absolute;top:909px;left:128px;white-space:nowrap" class="ft06">envolvidas na fabricação;de;medicamentos, incluindo os medicamentos experimentais.;</p>
<p style="position:absolute;top:946px;left:412px;white-space:nowrap" class="ft07"><b>Seção III;</b></p>
<p style="position:absolute;top:981px;left:392px;white-space:nowrap" class="ft07"><b>Das definições;</b></p>
<p style="position:absolute;top:1016px;left:170px;white-space:nowrap" class="ft06">Art. 3º Para fins desta Resolução e das instruções normativas vinculadas a ela,;</p>
<p style="position:absolute;top:1037px;left:128px;white-space:nowrap" class="ft06">aplicam-se;as seguintes definições:;</p>
<p style="position:absolute;top:1073px;left:170px;white-space:nowrap" class="ft06">I;-;acordo;técnico:;documento;que;define;responsabilidades,;atribuições,;direitos;e;</p>
<p style="position:absolute;top:1094px;left:128px;white-space:nowrap" class="ft06">deveres de/entre;contratante e;contratado em relação às atividades terceirizadas </p>
</div>
</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHtml)
	// links, err := parser.GetFontMap(r)
	// if err != nil {
	// 	panic(err)
	// }
	// for k, v := range *links {
	// 	if strings.TrimSpace(k) == CHAPTER_FONT {
	// 		font = v
	// 	}
	// }
	parser.ParseToFont(r)
}
