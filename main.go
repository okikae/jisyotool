/* vim:ts=4:
 * Author: 奈幾乃(uakms)
 * Created: 2015-04-09
 * Revised: 2024-12-13
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	reComment = regexp.MustCompile("^;.*|^$|^[\t\n\f\r ]+$") //コメント行
	reNote    = regexp.MustCompile("[\t\n\f\r ]+;.*")        //備考
	reSepa    = regexp.MustCompile(" /")                     //辞書の区切り文字
)

func simpleCreateJsonDict(jisyoPath string, flag string) []string {
	jsonDictArr := []string{}
	ifile, err := os.Open(jisyoPath)
	if err != nil {
		fmt.Println("ファイルを読み込めませんでした")
		os.Exit(0)
	}
	defer ifile.Close()
	scanner := bufio.NewScanner(ifile)

	for scanner.Scan() {
		commentedLine := reComment.MatchString(scanner.Text())
		if commentedLine == false {
			mbody := reNote.ReplaceAllString(scanner.Text(), "")
			supl := reNote.FindStringSubmatch(scanner.Text())
			suplstr := ""
			if len(supl) != 0 {
				reSplmnt := regexp.MustCompile("\\[(.*)\\]")
				supl2 := reSplmnt.FindAllStringSubmatch(supl[0], -1)
				if len(supl2) != 0 {
					reSubSep := regexp.MustCompile("\\|")
					suplstr = "[\"" +
						reSubSep.ReplaceAllString(supl2[0][1], "\", \"") + "\"]"
				} else {
					suplstr = "[]"
				}
			} else {
				suplstr = "[]"
			}
			pair := reSepa.Split(mbody, 2)
			if flag == "normal" {
				elem :=
					"[\"" + pair[0] + "\", \"" + pair[1] + "\", " + suplstr + "]"
				jsonDictArr = append(jsonDictArr, elem)
			} else if flag == "reverse" {
				elem :=
					"[\"" + pair[1] + "\", \"" + pair[0] + "\", " + suplstr + "]"
				jsonDictArr = append(jsonDictArr, elem)
			}
		}
	}
	return jsonDictArr
}

func simpleOutputJsonDict(jsonDictArr []string) {
	fmt.Println("[")
	for i := 0; i < len(jsonDictArr)-1; i++ {
		fmt.Println("  "+jsonDictArr[i]+",")
	}
	fmt.Println("  "+jsonDictArr[len(jsonDictArr)-1]+"\n]")
	return
}


func createJsonDict(jisyoPath string, flag string) []string {
	jsonDictArr := []string{}
	ifile, err := os.Open(jisyoPath)
	if err != nil {
		fmt.Println("ファイルを読み込めませんでした")
		os.Exit(0)
	}
	defer ifile.Close()
	scanner := bufio.NewScanner(ifile)

	for scanner.Scan() {
		var isExcludeLine bool = false
		if flag == "normal" {
			isExcludeLine = (strings.Contains(scanner.Text(), "新旧変換で不要？") || strings.Contains(scanner.Text(), "新旧誤変換"))
		} else if flag == "reverse" {
			isExcludeLine = (strings.Contains(scanner.Text(), "旧新誤変換") || strings.Contains(scanner.Text(), "旧新変換で不要？"))
		}
		commentedLine := reComment.MatchString(scanner.Text())
		if (commentedLine == false) && (isExcludeLine == false) {
			mbody := reNote.ReplaceAllString(scanner.Text(), "")
			supl := reNote.FindStringSubmatch(scanner.Text())
			suplstr := ""
			if len(supl) != 0 {
				reSplmnt := regexp.MustCompile("\\[(.*)\\]")
				supl2 := reSplmnt.FindAllStringSubmatch(supl[0], -1)
				if len(supl2) != 0 {
					reSubSep := regexp.MustCompile("\\|")
					suplstr = "[\"" +
						reSubSep.ReplaceAllString(supl2[0][1], "\", \"") + "\"]"
				} else {
					suplstr = "[]"
				}
			} else {
				suplstr = "[]"
			}
			pair := reSepa.Split(mbody, 2)
			if flag == "normal" {
				elem :=
					"[\"" + pair[0] + "\", \"" + pair[1] + "\", " + suplstr + "]"
				jsonDictArr = append(jsonDictArr, elem)
			} else if flag == "reverse" {
				elem :=
					"[\"" + pair[1] + "\", \"" + pair[0] + "\", " + suplstr + "]"
				jsonDictArr = append(jsonDictArr, elem)
			}
		}
	}
	return jsonDictArr
}

func outputJsonDict(jsonDictArr []string, pref string, fname string) {
	ofile, err := os.Create(fname)
	if err != nil {
		fmt.Println("ファイルを書き込めませんでした")
		os.Exit(0)
	}
	defer ofile.Close()
	writer := bufio.NewWriter(ofile)
	fmt.Fprint(writer, pref+" =\n[\n")
	for i := 0; i < len(jsonDictArr)-1; i++ {
		fmt.Fprint(writer, "  "+jsonDictArr[i]+",\n")
	}
	fmt.Fprint(writer, "  "+jsonDictArr[len(jsonDictArr)-1]+"\n]\n")
	writer.Flush()
	return
}

func checkerDict(jisyoPath string) {
	var (
		extractArr   = []string{}
		reIncomplete = regexp.MustCompile("　|；| /.\\S*;|^.\\S*/|\\s$")
	)

	ifile, err := os.Open(jisyoPath)
	if err != nil {
		fmt.Println("ファイルを読み込めませんでした")
		os.Exit(0)
	}
	defer ifile.Close()
	scanner := bufio.NewScanner(ifile)

	for scanner.Scan() {
		fault := reIncomplete.MatchString(scanner.Text())
		if fault == true {
			extractArr = append(extractArr, scanner.Text())
		}
	}

	if len(extractArr) == 0 {
		fmt.Println("(^_^) 辞書の書式に問題はなさそうです")
	} else {
		fmt.Println("(>_<) 次の行に問題がありそうです")
		for _, element := range extractArr {
			fmt.Printf("%s\n", element)
		}
	}
	return
}

func checkDictDuplicate(mainDictArr [][]string) {
	duplicatedarr := [][]string{}
	for i := 0; i < len(mainDictArr); i++ {
		if mainDictArr[i][0] == mainDictArr[i][1] {
			duplicatedarr = append(duplicatedarr, mainDictArr[i])
		}
		for j := i + 1; j < len(mainDictArr); j++ {
			if len(mainDictArr[i]) != 1 {
				if (mainDictArr[i][0] == mainDictArr[j][0]) &&
					(mainDictArr[i][1] == mainDictArr[j][1]) {
					duplicatedarr = append(duplicatedarr, mainDictArr[i])
				}
			}
		}
	}
	if len(duplicatedarr) == 0 {
		fmt.Println("\n(^_^) 辞書に重複はなさそうです")
	} else {
		fmt.Println("\n(>_<) 次の単語に重複登録、もしくは key と value が同じ状態のものがありそうです")
		for _, dups := range duplicatedarr {
			fmt.Printf("%v\n", dups[0])
		}
	}
	return
}

func creatDict(jisyoPath string) [][]string {
	mainDictArr := [][]string{}
	ifile, err := os.Open(jisyoPath)
	if err != nil {
		fmt.Println("ファイルを読み込めませんでした")
		os.Exit(0)
	}
	defer ifile.Close()
	scanner := bufio.NewScanner(ifile)

	for scanner.Scan() {
		commentedLine := reComment.MatchString(scanner.Text())
		if commentedLine == false {
			str := reNote.ReplaceAllString(scanner.Text(), "")
			pair := reSepa.Split(str, 2)
			mainDictArr = append(mainDictArr, pair)
		}
	}
	return mainDictArr
}

func outputDict(mainDictArr [][]string, flag string) {
	for _, element := range mainDictArr {
		if flag == "normal" {
			fmt.Printf("%s %s\n", element[0], element[1])
		} else if flag == "reverse" {
			fmt.Printf("%s %s\n", element[1], element[0])
		}
	}
	return
}

func printUsage() {
	fmt.Println("Usage: dictutil option inputfile")
	fmt.Println("  option:")
	fmt.Println("          -jn [normal]  カラム1 カラム2 の順で JSON 出力")
	fmt.Println("          -jr [reverse] カラム2 カラム1 の順で JSON 出力")
	fmt.Println("          -n [normal]   カラム1 カラム2 の順で出力")
	fmt.Println("          -r [reverse]  カラム2 カラム1 の順で出力")
	fmt.Println("          -j [json]     JSON形式でファイルに出力")
	fmt.Println("          -j2 [json]    JSON形式でファイルに出力(分離)")
	fmt.Println("          -t [json]     JSON形式でts拡張子のファイルに出力")
	fmt.Println("          -t2 [json]    JSON形式でts拡張子のファイルに出力(分離)")
	fmt.Println("          -c [check]    辞書をチェック")
	fmt.Println("          -l [length]   要素数を出力")
	return
}

func parseArgument() {
	if len(os.Args) <= 2 || len(os.Args) >= 4 {
		printUsage()
		return
	}

	if len(os.Args) == 3 {
		jisyoPath := os.Args[2]
		jisyoArr := strings.Split(jisyoPath, "/")
		jisyoName := jisyoArr[len(jisyoArr)-1]
		mainArr := creatDict(jisyoPath)

		switch {
		case os.Args[1] == "-n":
			outputDict(mainArr, "normal")
			return

		case os.Args[1] == "-r":
			outputDict(mainArr, "reverse")
			return

		case os.Args[1] == "-l":
			fmt.Println(len(mainArr))
			return

		case os.Args[1] == "-c":
			checkerDict(jisyoPath)
			checkDictDuplicate(mainArr)
			return

		case os.Args[1] == "-jn":
			jsonArr := simpleCreateJsonDict(jisyoPath, "normal")
			simpleOutputJsonDict(jsonArr)
			return

		case os.Args[1] == "-jr":
			jsonArr := simpleCreateJsonDict(jisyoPath, "reverse")
			simpleOutputJsonDict(jsonArr)
			return

		case os.Args[1] == "-j":
			re := regexp.MustCompile("(.*)-jisyo")
			matched := re.FindStringSubmatch(jisyoName)
			jsonArr := createJsonDict(jisyoPath, "normal")
			if matched[1] == "kana" {
				outputJsonDict(jsonArr, "var kanaArray", "dic-kana.js")
			} else if matched[1] == "kanji" {
				outputJsonDict(jsonArr, "var kanjiArray", "dic-kanji.js")
			} else {
				fmt.Println("JSON出力に対応している辞書は")
				fmt.Println("[kana, kanji] のみです。")
			}
			return

		case os.Args[1] == "-j2":
			re := regexp.MustCompile("(.*)-jisyo")
			matched := re.FindStringSubmatch(jisyoName)
			if matched[1] == "kana" {
				jsonArr0 := createJsonDict(jisyoPath, "reverse")
				outputJsonDict(jsonArr0,
					"var toModernKanaArray", "dic-to-modern-kana.js")
				jsonArr1 := createJsonDict(jisyoPath, "normal")
				outputJsonDict(jsonArr1,
					"var toTradKanaArray", "dic-to-trad-kana.js")
			} else if matched[1] == "kanji" {
				jsonArr0 := createJsonDict(jisyoPath, "reverse")
				outputJsonDict(jsonArr0,
					"var toNewKanjiArray", "dic-to-new-kanji.js")
				jsonArr1 := createJsonDict(jisyoPath, "normal")
				outputJsonDict(jsonArr1,
					"var toOldKanjiArray", "dic-to-old-kanji.js")
			} else {
				fmt.Println("JSON出力に対応している辞書は")
				fmt.Println("[kana, kanji] のみです。")
			}
			return

		case os.Args[1] == "-t":
			re := regexp.MustCompile("(.*)-jisyo")
			matched := re.FindStringSubmatch(jisyoName)
			jsonArr := createJsonDict(jisyoPath, "normal")
			if matched[1] == "kana" {
				outputJsonDict(jsonArr, "export const kanaArray: Array<[string, string, Array<string>]>", "kanajisyo.ts")
			} else if matched[1] == "kanji" {
				outputJsonDict(jsonArr, "export const kanjiArray: Array<[string, string, Array<string>]>", "kanjijisyo.ts")
			} else {
				fmt.Println("JSON出力に対応している辞書は")
				fmt.Println("[kana, kanji] のみです。")
			}
			return

		case os.Args[1] == "-t2":
			re := regexp.MustCompile("(.*)-jisyo")
			matched := re.FindStringSubmatch(jisyoName)
			if matched[1] == "kana" {
				jsonArr0 := createJsonDict(jisyoPath, "reverse")
				outputJsonDict(jsonArr0,
					"export const toModernKanaArray: Array<[string, string, Array<string>]>",
					"modernkanajisyo.ts")
				jsonArr1 := createJsonDict(jisyoPath, "normal")
				outputJsonDict(jsonArr1,
					"export const toTradKanaArray: Array<[string, string, Array<string>]>",
					"tradkanajisyo.ts")
			} else if matched[1] == "kanji" {
				jsonArr0 := createJsonDict(jisyoPath, "reverse")
				outputJsonDict(jsonArr0,
					"export const toNewKanjiArray: Array<[string, string, Array<string>]>",
					"newkanjijisyo.ts")
				jsonArr1 := createJsonDict(jisyoPath, "normal")
				outputJsonDict(jsonArr1,
					"export const toOldKanjiArray: Array<[string, string, Array<string>]>",
					"oldkanjijisyo.ts")
			} else {
				fmt.Println("JSON出力に対応している辞書は")
				fmt.Println("[kana, kanji] のみです。")
			}
			return
		}
		printUsage()
		return
	}
}

func main() {
	parseArgument()
	return
}
