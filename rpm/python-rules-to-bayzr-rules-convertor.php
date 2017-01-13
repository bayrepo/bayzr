<?php

function convMsm($val, $str){
    $mult = 1;
    if ($str=="mn") $mult = 1;
    if ($str=="h") $mult = 60;
    $res = $val * $mult;
    return $res . "min";
}

function conSev($val){
    switch($val[0]){
	case 'C':
	    return "INFO";
	case 'I':
	    return "INFO";
	case 'R':
	    return "INFO";
	case 'E':
	    return "MAJOR";
	case 'W':
	    return "MINOR";
	case 'F':
	    return "CRITICAL";
    }
    return "MINOR";
}

function conTp($val){
    switch($val[0]){
	case 'C':
	    return "CODE_SMELL";
	case 'I':
	    return "CODE_SMELL";
	case 'R':
	    return "CODE_SMELL";
	case 'E':
	    return "BUG";
	case 'W':
	    return "CODE_SMELL";
	case 'F':
	    return "BUG";
    }
    return "CODE_SMELL";
}

$xml = simplexml_load_file("python-model2.xml");
$counter = 0;
$arr=array();
echo "<rules>\n";
foreach ($xml->chc as $chc_parent) {
    foreach ($chc_parent->chc as $chc_2) {
	foreach ($chc_2->chc as $chc_3) {
	    if (in_array(strtolower($chc_3->rule_key), $arr)){
		continue;
	    }
	    $arr[]=strtolower($chc_3->rule_key);
	    echo "  <rule>\n";
	    echo "     <key>" . $chc_3->rule_key . "</key>\n";
	    $file = @file_get_contents('http://pylint-messages.wikidot.com/messages:'.strtolower($chc_3->rule_key));
	    if ($file!=""){
		$doc = new DOMDocument();
		$doc->loadHTML($file);
		$els=$doc->getElementsByTagName('code');
		$name = $chc_3->rule_key;
		foreach ($els as $el) {
		    $name = trim(str_replace(array("\n", "\r"), '', $el->nodeValue));
		    echo "     <name><![CDATA[".$name."]]></name>\n";
		}
		$elm = $doc->getElementById('page-content');
		if ($elm!=NULL){
			$els=$elm->getElementsByTagName('p');
			$result = "";
			foreach ($els as $el) {
			    $result = $result." ".$el->nodeValue;
			}
			$rs = trim(str_replace(array("\n", "\r"), ' ', $result));
			echo "     <description><![CDATA[".$rs."]]></description>\n";
		} else {
		    echo "     <description><![CDATA[".$name."]]></description>\n";
		}
	    } else {
		echo "     <name><![CDATA[".$chc_3->rule_key."]]></name>\n";
		echo "     <description><![CDATA[".$chc_3->rule_key."]]></description>\n";
	    }
	    echo "     <internalKey>" . $chc_3->rule_key . "</internalKey>\n";
	    echo "     <remediationFunction>" . strtoupper($chc_3->prop[0]->txt) . "</remediationFunction>\n";
	    echo "     <remediationFunctionGapMultiplier>" . convMsm($chc_3->prop[1]->val , $chc_3->prop[1]->txt). "</remediationFunctionGapMultiplier>\n";
	    echo "     <severity>" . conSev($chc_3->rule_key). "</severity>\n";
	    echo "     <type>" . conTp($chc_3->rule_key). "</type>\n";
	    echo "     <tag>bayzr</tag>\n";
	    echo "     <tag>python</tag>\n";
	    echo "  </rule>\n";
	    $counter++;
       }
    }
}
echo "</rules>\n";
?>